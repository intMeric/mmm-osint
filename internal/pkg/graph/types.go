package graph

import (
	"fmt"
	"strings"
)

type NodeType string

const (
	NodeTypeURL  NodeType = "URL"
	NodeTypeUser NodeType = "User"
)

func (nt NodeType) String() string {
	return string(nt)
}

func (nt NodeType) IsValid() bool {
	switch nt {
	case NodeTypeURL, NodeTypeUser:
		return true
	default:
		return false
	}
}

type Node struct {
	Type        NodeType `json:"type"`
	DisplayName string   `json:"displayName"`
	ID          string   `json:"id"`
	Location    string   `json:"location,omitempty"`
}

func (n *Node) Validate() error {
	if !n.Type.IsValid() {
		return fmt.Errorf("invalid node type: %s", n.Type)
	}

	if strings.TrimSpace(n.DisplayName) == "" {
		return fmt.Errorf("displayName cannot be empty")
	}

	if strings.TrimSpace(n.ID) == "" {
		return fmt.Errorf("id cannot be empty")
	}

	return nil
}

type RelationType string

const (
	RelationTypeConnectedTo RelationType = "CONNECTED_TO"
	RelationTypeRelatesTo   RelationType = "RELATES_TO"
	RelationTypeLinkedTo    RelationType = "LINKED_TO"
)

func (rt RelationType) String() string {
	return string(rt)
}

func (rt RelationType) IsValid() bool {
	switch rt {
	case RelationTypeConnectedTo, RelationTypeRelatesTo, RelationTypeLinkedTo:
		return true
	default:
		return false
	}
}

type Relation struct {
	Type     RelationType `json:"type"`
	SourceID string       `json:"sourceId"`
	TargetID string       `json:"targetId"`
}

func (r *Relation) Validate() error {
	if !r.Type.IsValid() {
		return fmt.Errorf("invalid relation type: %s", r.Type)
	}

	if strings.TrimSpace(r.SourceID) == "" {
		return fmt.Errorf("sourceId cannot be empty")
	}

	if strings.TrimSpace(r.TargetID) == "" {
		return fmt.Errorf("targetId cannot be empty")
	}

	if r.SourceID == r.TargetID {
		return fmt.Errorf("sourceId and targetId cannot be the same")
	}

	return nil
}