package graph_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"mmm-osint/internal/pkg/graph"
)

var _ = Describe("Types", func() {
	Describe("NodeType", func() {
		Context("with valid node types", func() {
			It("should validate URL type", func() {
				nodeType := graph.NodeTypeURL
				Expect(nodeType.IsValid()).To(BeTrue())
				Expect(nodeType.String()).To(Equal("URL"))
			})

			It("should validate User type", func() {
				nodeType := graph.NodeTypeUser
				Expect(nodeType.IsValid()).To(BeTrue())
				Expect(nodeType.String()).To(Equal("User"))
			})
		})

		Context("with invalid node types", func() {
			It("should reject invalid types", func() {
				nodeType := graph.NodeType("INVALID")
				Expect(nodeType.IsValid()).To(BeFalse())
			})

			It("should reject empty types", func() {
				nodeType := graph.NodeType("")
				Expect(nodeType.IsValid()).To(BeFalse())
			})
		})
	})

	Describe("Node", func() {
		var node *graph.Node

		BeforeEach(func() {
			node = &graph.Node{
				Type:        graph.NodeTypeURL,
				DisplayName: "Test URL",
				ID:          "test-id-123",
				Location:    "db.collection",
			}
		})

		Context("with valid node", func() {
			It("should validate successfully", func() {
				err := node.Validate()
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("with invalid node type", func() {
			BeforeEach(func() {
				node.Type = "INVALID"
			})

			It("should return validation error", func() {
				err := node.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("invalid node type"))
			})
		})

		Context("with empty displayName", func() {
			It("should return validation error for empty string", func() {
				node.DisplayName = ""
				err := node.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("displayName cannot be empty"))
			})

			It("should return validation error for whitespace-only string", func() {
				node.DisplayName = "   "
				err := node.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("displayName cannot be empty"))
			})
		})

		Context("with empty ID", func() {
			It("should return validation error for empty string", func() {
				node.ID = ""
				err := node.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("id cannot be empty"))
			})

			It("should return validation error for whitespace-only string", func() {
				node.ID = "   "
				err := node.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("id cannot be empty"))
			})
		})

		Context("with empty location", func() {
			It("should validate successfully (location is optional)", func() {
				node.Location = ""
				err := node.Validate()
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Describe("RelationType", func() {
		Context("with valid relation types", func() {
			It("should validate CONNECTED_TO type", func() {
				relationType := graph.RelationTypeConnectedTo
				Expect(relationType.IsValid()).To(BeTrue())
				Expect(relationType.String()).To(Equal("CONNECTED_TO"))
			})

			It("should validate RELATES_TO type", func() {
				relationType := graph.RelationTypeRelatesTo
				Expect(relationType.IsValid()).To(BeTrue())
				Expect(relationType.String()).To(Equal("RELATES_TO"))
			})

			It("should validate LINKED_TO type", func() {
				relationType := graph.RelationTypeLinkedTo
				Expect(relationType.IsValid()).To(BeTrue())
				Expect(relationType.String()).To(Equal("LINKED_TO"))
			})
		})

		Context("with invalid relation types", func() {
			It("should reject invalid types", func() {
				relationType := graph.RelationType("INVALID")
				Expect(relationType.IsValid()).To(BeFalse())
			})

			It("should reject empty types", func() {
				relationType := graph.RelationType("")
				Expect(relationType.IsValid()).To(BeFalse())
			})
		})
	})

	Describe("Relation", func() {
		var relation *graph.Relation

		BeforeEach(func() {
			relation = &graph.Relation{
				Type:     graph.RelationTypeConnectedTo,
				SourceID: "source-id-123",
				TargetID: "target-id-456",
			}
		})

		Context("with valid relation", func() {
			It("should validate successfully", func() {
				err := relation.Validate()
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("with invalid relation type", func() {
			BeforeEach(func() {
				relation.Type = "INVALID"
			})

			It("should return validation error", func() {
				err := relation.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("invalid relation type"))
			})
		})

		Context("with empty sourceId", func() {
			It("should return validation error for empty string", func() {
				relation.SourceID = ""
				err := relation.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("sourceId cannot be empty"))
			})

			It("should return validation error for whitespace-only string", func() {
				relation.SourceID = "   "
				err := relation.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("sourceId cannot be empty"))
			})
		})

		Context("with empty targetId", func() {
			It("should return validation error for empty string", func() {
				relation.TargetID = ""
				err := relation.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("targetId cannot be empty"))
			})

			It("should return validation error for whitespace-only string", func() {
				relation.TargetID = "   "
				err := relation.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("targetId cannot be empty"))
			})
		})

		Context("with same sourceId and targetId", func() {
			BeforeEach(func() {
				relation.TargetID = relation.SourceID
			})

			It("should return validation error", func() {
				err := relation.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("sourceId and targetId cannot be the same"))
			})
		})
	})
})