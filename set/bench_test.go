package set

import (
	"fmt"
	"testing"

	"github.com/jonwraymond/toolfoundation/adapter"
)

func BenchmarkToolset_Add(b *testing.B) {
	ts := New("bench")
	tool := &adapter.CanonicalTool{Name: "test", Namespace: "ns"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ts.Add(tool)
	}
}

func BenchmarkToolset_Get(b *testing.B) {
	ts := New("bench")
	for i := 0; i < 1000; i++ {
		ts.Add(&adapter.CanonicalTool{
			Name:      fmt.Sprintf("tool%d", i),
			Namespace: "ns",
		})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ts.Get("ns:tool500")
	}
}

func BenchmarkToolset_Tools(b *testing.B) {
	ts := New("bench")
	for i := 0; i < 1000; i++ {
		ts.Add(&adapter.CanonicalTool{
			Name:      fmt.Sprintf("tool%d", i),
			Namespace: "ns",
		})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ts.Tools()
	}
}

func BenchmarkToolset_Filter(b *testing.B) {
	ts := New("bench")
	for i := 0; i < 1000; i++ {
		tags := []string{"common"}
		if i%2 == 0 {
			tags = append(tags, "even")
		}
		ts.Add(&adapter.CanonicalTool{
			Name:      fmt.Sprintf("tool%d", i),
			Namespace: "ns",
			Tags:      tags,
		})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ts.Filter(TagsAny("even"))
	}
}

func BenchmarkBuilder_Build(b *testing.B) {
	for _, size := range []int{100, 500, 1000} {
		tools := make([]*adapter.CanonicalTool, size)
		for i := range tools {
			tools[i] = &adapter.CanonicalTool{
				Name:      fmt.Sprintf("tool%d", i),
				Namespace: "ns",
				Tags:      []string{"tag1", "tag2"},
			}
		}

		b.Run(fmt.Sprintf("tools_%d", size), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = NewBuilder("bench").
					FromTools(tools).
					Build()
			}
		})
	}
}

func BenchmarkBuilder_Build_WithFilters(b *testing.B) {
	tools := make([]*adapter.CanonicalTool, 1000)
	for i := range tools {
		ns := "ns1"
		if i%2 == 0 {
			ns = "ns2"
		}
		tools[i] = &adapter.CanonicalTool{
			Name:      fmt.Sprintf("tool%d", i),
			Namespace: ns,
			Tags:      []string{"tag1", "tag2"},
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = NewBuilder("bench").
			FromTools(tools).
			WithNamespace("ns1").
			WithTags([]string{"tag1"}).
			Build()
	}
}

func BenchmarkBuilder_Build_WithPolicy(b *testing.B) {
	tools := make([]*adapter.CanonicalTool, 1000)
	for i := range tools {
		tags := []string{"safe"}
		if i%10 == 0 {
			tags = append(tags, "dangerous")
		}
		tools[i] = &adapter.CanonicalTool{
			Name:      fmt.Sprintf("tool%d", i),
			Namespace: "ns",
			Tags:      tags,
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = NewBuilder("bench").
			FromTools(tools).
			WithPolicy(DenyTags("dangerous")).
			Build()
	}
}

func BenchmarkNamespaceFilter(b *testing.B) {
	tool := &adapter.CanonicalTool{Name: "test", Namespace: "github"}
	filter := NamespaceFilter("github", "gitlab", "bitbucket")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filter(tool)
	}
}

func BenchmarkTagsAll(b *testing.B) {
	tool := &adapter.CanonicalTool{
		Name: "test",
		Tags: []string{"read", "write", "admin", "api"},
	}
	filter := TagsAll("read", "write")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filter(tool)
	}
}
