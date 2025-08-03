package graph

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type DriverWithContext interface {
	VerifyConnectivity(ctx context.Context) error
	NewSession(ctx context.Context, config neo4j.SessionConfig) neo4j.SessionWithContext
	Close(ctx context.Context) error
}

type Neo4jGraph struct {
	driver DriverWithContext
	mu     sync.RWMutex
	config *Neo4jConfig
}

type Neo4jConfig struct {
	URI      string
	Username string
	Password string
	Database string
}

func NewNeo4jGraph(config *Neo4jConfig) (*Neo4jGraph, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	auth := neo4j.BasicAuth(config.Username, config.Password, "")
	
	driver, err := neo4j.NewDriverWithContext(
		config.URI,
		auth,
		func(config *neo4j.Config) {
			config.MaxConnectionLifetime = 5 * time.Minute
			config.MaxConnectionPoolSize = 50
			config.ConnectionAcquisitionTimeout = 2 * time.Minute
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create Neo4j driver: %w", err)
	}

	graph := &Neo4jGraph{
		driver: driver,
		config: config,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := graph.verifyConnectivity(ctx); err != nil {
		driver.Close(ctx)
		return nil, fmt.Errorf("failed to verify connectivity: %w", err)
	}

	return graph, nil
}

func (g *Neo4jGraph) verifyConnectivity(ctx context.Context) error {
	g.mu.RLock()
	driver := g.driver
	g.mu.RUnlock()

	return driver.VerifyConnectivity(ctx)
}

func (g *Neo4jGraph) ensureConnection(ctx context.Context) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if err := g.driver.VerifyConnectivity(ctx); err != nil {
		g.driver.Close(ctx)

		auth := neo4j.BasicAuth(g.config.Username, g.config.Password, "")
		
		driver, newErr := neo4j.NewDriverWithContext(
			g.config.URI,
			auth,
			func(config *neo4j.Config) {
				config.MaxConnectionLifetime = 5 * time.Minute
				config.MaxConnectionPoolSize = 50
				config.ConnectionAcquisitionTimeout = 2 * time.Minute
			},
		)
		if newErr != nil {
			return fmt.Errorf("failed to reconnect to Neo4j: %w", newErr)
		}

		g.driver = driver
		
		if verifyErr := g.driver.VerifyConnectivity(ctx); verifyErr != nil {
			return fmt.Errorf("failed to verify reconnection: %w", verifyErr)
		}
	}

	return nil
}

func (g *Neo4jGraph) CreateNode(ctx context.Context, node *Node) error {
	if node == nil {
		return fmt.Errorf("node cannot be nil")
	}

	if err := node.Validate(); err != nil {
		return fmt.Errorf("invalid node: %w", err)
	}

	if err := g.ensureConnection(ctx); err != nil {
		return fmt.Errorf("connection error: %w", err)
	}

	g.mu.RLock()
	driver := g.driver
	g.mu.RUnlock()

	session := driver.NewSession(ctx, neo4j.SessionConfig{
		DatabaseName: g.config.Database,
	})
	defer session.Close(ctx)

	query := fmt.Sprintf(`
		MERGE (n:%s {id: $id})
		SET n.displayName = $displayName,
		    n.location = $location,
		    n.updated_at = datetime()
		ON CREATE SET n.created_at = datetime()
		RETURN n
	`, node.Type)

	parameters := map[string]any{
		"id":          node.ID,
		"displayName": node.DisplayName,
		"location":    node.Location,
	}

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, query, parameters)
		if err != nil {
			return nil, err
		}

		return result.Consume(ctx)
	})

	if err != nil {
		return fmt.Errorf("failed to create node: %w", err)
	}

	return nil
}

func (g *Neo4jGraph) CreateRelation(ctx context.Context, relation *Relation) error {
	if relation == nil {
		return fmt.Errorf("relation cannot be nil")
	}

	if err := relation.Validate(); err != nil {
		return fmt.Errorf("invalid relation: %w", err)
	}

	if err := g.ensureConnection(ctx); err != nil {
		return fmt.Errorf("connection error: %w", err)
	}

	g.mu.RLock()
	driver := g.driver
	g.mu.RUnlock()

	session := driver.NewSession(ctx, neo4j.SessionConfig{
		DatabaseName: g.config.Database,
	})
	defer session.Close(ctx)

	query := fmt.Sprintf(`
		MATCH (source {id: $sourceId})
		MATCH (target {id: $targetId})
		MERGE (source)-[r:%s]->(target)
		ON CREATE SET r.created_at = datetime()
		ON MATCH SET r.updated_at = datetime()
		RETURN r
	`, relation.Type)

	parameters := map[string]any{
		"sourceId": relation.SourceID,
		"targetId": relation.TargetID,
	}

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, query, parameters)
		if err != nil {
			return nil, err
		}

		return result.Consume(ctx)
	})

	if err != nil {
		return fmt.Errorf("failed to create relation: %w", err)
	}

	return nil
}

func (g *Neo4jGraph) GetNode(ctx context.Context, id string) (*Node, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("id cannot be empty")
	}

	if err := g.ensureConnection(ctx); err != nil {
		return nil, fmt.Errorf("connection error: %w", err)
	}

	g.mu.RLock()
	driver := g.driver
	g.mu.RUnlock()

	session := driver.NewSession(ctx, neo4j.SessionConfig{
		DatabaseName: g.config.Database,
	})
	defer session.Close(ctx)

	query := `
		MATCH (n {id: $id})
		RETURN n.id as id, n.displayName as displayName, n.location as location, labels(n) as labels
	`

	parameters := map[string]any{
		"id": id,
	}

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, query, parameters)
		if err != nil {
			return nil, err
		}

		record, err := result.Single(ctx)
		if err != nil {
			return nil, err
		}

		return record, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get node: %w", err)
	}

	record := result.(*neo4j.Record)
	
	nodeID, _ := record.Get("id")
	displayName, _ := record.Get("displayName")
	location, _ := record.Get("location")
	labels, _ := record.Get("labels")

	labelsList := labels.([]any)
	if len(labelsList) == 0 {
		return nil, fmt.Errorf("node has no labels")
	}

	nodeType := NodeType(labelsList[0].(string))
	
	node := &Node{
		Type:        nodeType,
		DisplayName: displayName.(string),
		ID:          nodeID.(string),
		Location:    location.(string),
	}

	return node, nil
}

func (g *Neo4jGraph) NodeExists(ctx context.Context, id string) (bool, error) {
	if strings.TrimSpace(id) == "" {
		return false, fmt.Errorf("id cannot be empty")
	}

	if err := g.ensureConnection(ctx); err != nil {
		return false, fmt.Errorf("connection error: %w", err)
	}

	g.mu.RLock()
	driver := g.driver
	g.mu.RUnlock()

	session := driver.NewSession(ctx, neo4j.SessionConfig{
		DatabaseName: g.config.Database,
	})
	defer session.Close(ctx)

	query := `
		MATCH (n {id: $id})
		RETURN count(n) > 0 as exists
	`

	parameters := map[string]any{
		"id": id,
	}

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, query, parameters)
		if err != nil {
			return nil, err
		}

		record, err := result.Single(ctx)
		if err != nil {
			return nil, err
		}

		exists, _ := record.Get("exists")
		return exists.(bool), nil
	})

	if err != nil {
		return false, fmt.Errorf("failed to check node existence: %w", err)
	}

	return result.(bool), nil
}

func (g *Neo4jGraph) Close(ctx context.Context) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.driver != nil {
		return g.driver.Close(ctx)
	}
	return nil
}