package graph_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"mmm-osint/internal/pkg/graph"
)

var _ = Describe("Factory", func() {

	Describe("NewGraph", func() {
		Context("with Neo4j graph type", func() {
			BeforeEach(func() {
				os.Setenv("NEO4J_URI", "neo4j://localhost:7687")
				os.Setenv("NEO4J_USERNAME", "neo4j")
				os.Setenv("NEO4J_PASSWORD", "testpassword")
				os.Setenv("NEO4J_DATABASE", "neo4j")
			})

			AfterEach(func() {
				os.Unsetenv("NEO4J_URI")
				os.Unsetenv("NEO4J_USERNAME")
				os.Unsetenv("NEO4J_PASSWORD")
				os.Unsetenv("NEO4J_DATABASE")
			})

			It("should attempt to create a Neo4j graph instance", func() {
				g, err := graph.NewGraph(graph.Neo4jGraphType)
				
				Expect(g).To(BeNil())
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed to verify connectivity"))
			})
		})

		Context("with unsupported graph type", func() {
			It("should return an error", func() {
				g, err := graph.NewGraph("unsupported")
				
				Expect(g).To(BeNil())
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("unsupported graph type"))
			})
		})
	})

	Describe("NewGraphWithConfig", func() {
		Context("with Neo4j graph type and valid config", func() {
			It("should attempt to create a Neo4j graph instance", func() {
				config := &graph.Neo4jConfig{
					URI:      "neo4j://localhost:7687",
					Username: "neo4j",
					Password: "testpassword",
					Database: "neo4j",
				}
				
				g, err := graph.NewGraphWithConfig(graph.Neo4jGraphType, config)
				
				Expect(g).To(BeNil())
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed to verify connectivity"))
			})
		})

		Context("with Neo4j graph type and invalid config", func() {
			It("should return an error", func() {
				g, err := graph.NewGraphWithConfig(graph.Neo4jGraphType, "invalid")
				
				Expect(g).To(BeNil())
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("invalid config type"))
			})
		})

		Context("with unsupported graph type", func() {
			It("should return an error", func() {
				g, err := graph.NewGraphWithConfig("unsupported", nil)
				
				Expect(g).To(BeNil())
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("unsupported graph type"))
			})
		})
	})
})