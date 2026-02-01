package gendoc_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/pseudomuto/protokit"
	"github.com/pseudomuto/protokit/utils"
	. "github.com/stillmatic/protoc-gen-doc"
	"github.com/stretchr/testify/require"
)

func TestRenderers(t *testing.T) {
	set, err := utils.LoadDescriptorSet("fixtures", "fileset.pb")
	require.NoError(t, err)

	os.Mkdir("./tmp", os.ModePerm)

	req := utils.CreateGenRequest(set, "Booking.proto", "Vehicle.proto")
	result := protokit.ParseCodeGenRequest(req)
	template := NewTemplate(result)

	for _, r := range []RenderType{
		RenderTypeDocBook,
		RenderTypeHTML,
		RenderTypeJSON,
		RenderTypeMarkdown,
	} {
		_, err := RenderTemplate(r, template, "")
		require.NoError(t, err)
	}
}

func TestNewRenderType(t *testing.T) {
	expected := []RenderType{
		RenderTypeDocBook,
		RenderTypeHTML,
		RenderTypeJSON,
		RenderTypeMarkdown,
	}

	supplied := []string{"docbook", "html", "json", "markdown"}

	for idx, input := range supplied {
		rt, err := NewRenderType(input)
		require.Nil(t, err)
		require.Equal(t, expected[idx], rt)
	}
}

func TestNewRenderTypeUnknown(t *testing.T) {
	rt, err := NewRenderType("/some/template.tmpl")
	require.Zero(t, rt)
	require.Error(t, err)
}

// TestMarkdownNoMultipleBlankLines verifies that the markdown template does not
// generate multiple consecutive blank lines (MD012 lint rule).
func TestMarkdownNoMultipleBlankLines(t *testing.T) {
	set, err := utils.LoadDescriptorSet("fixtures", "fileset.pb")
	require.NoError(t, err)

	req := utils.CreateGenRequest(set, "Booking.proto", "Vehicle.proto")
	result := protokit.ParseCodeGenRequest(req)
	template := NewTemplate(result)

	output, err := RenderTemplate(RenderTypeMarkdown, template, "")
	require.NoError(t, err)

	// Write output to file for debugging
	os.WriteFile("./tmp/test_output.md", output, 0644)

	// Check for 3+ consecutive newlines (which would mean 2+ blank lines)
	multipleBlankLines := regexp.MustCompile(`\n{3,}`)
	matches := multipleBlankLines.FindAllStringIndex(string(output), -1)

	if len(matches) > 0 {
		// Find the context around each match for better error messages
		outputStr := string(output)
		for i, match := range matches {
			start := match[0] - 50
			if start < 0 {
				start = 0
			}
			end := match[1] + 50
			if end > len(outputStr) {
				end = len(outputStr)
			}
			t.Errorf("Found multiple consecutive blank lines at position %d (match %d):\n...%s...",
				match[0], i+1, outputStr[start:end])
		}
	}
}
