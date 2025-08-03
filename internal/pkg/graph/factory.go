package graph

import (
	"fmt"
	"mmm-osint/internal/pkg/env"
)

type GraphType string

const (
	Neo4jGraphType GraphType = "neo4j"
)

func NewGraph(graphType GraphType) (Graph, error) {
	switch graphType {
	case Neo4jGraphType:
		config := &Neo4jConfig{
			URI:      env.GetOrDefault("NEO4J_URI", "neo4j://localhost:7687"),
			Username: env.GetOrDefault("NEO4J_USERNAME", "neo4j"),
			Password: env.GetOrDefault("NEO4J_PASSWORD", "password"),
			Database: env.GetOrDefault("NEO4J_DATABASE", "neo4j"),
		}
		return NewNeo4jGraph(config)
	default:
		return nil, fmt.Errorf("unsupported graph type: %s", graphType)
	}
}

func NewGraphWithConfig(graphType GraphType, config any) (Graph, error) {
	switch graphType {
	case Neo4jGraphType:
		neo4jConfig, ok := config.(*Neo4jConfig)
		if !ok {
			return nil, fmt.Errorf("invalid config type for Neo4j graph, expected *Neo4jConfig")
		}
		return NewNeo4jGraph(neo4jConfig)
	default:
		return nil, fmt.Errorf("unsupported graph type: %s", graphType)
	}
}