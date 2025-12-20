package api

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

//go:embed public/*
var staticFiles embed.FS

// ============================================================================
// DATA STRUCTURES
// ============================================================================
type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	PartNumber  string  `json:"part_number"`
	Price       float64 `json:"price"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
}

// ============================================================================
// HTML TEMPLATES
// ============================================================================

// Sidebar template - used across all pages for consistent navigation
const sidebarTemplate = `
<div class="sidebar">
    <p><a href="/">Alaska Hoffman</a></p>
    
    <form action="/search" method="GET">
        <input type="text" name="q" value="{{.Query}}">
    </form>
    
    <nav>
        <div><a href="/poetry">Poetry</a></div>
        <div><a href="/boma">Boma (2025)</a></div>
        <s><div><a href="/search?q=">Capsule 21 (2022)</a></div></s>
        <s><div><a href="/search?q=">Superchief Gallery (2023)</a></div></s>
        <s><div><a href="/search?q=">COEX, Korea (2023)</a></div></s>
        <s><div><a href="/search?q=">DX Singularity (2024)</a></div></s>
        <s><div><a href="/search?q=">DX Terminal (2025)</a></div></s>
        <s><div><a href="/search?q=">dark pressure rising (2025)</a></div></s>
        <br>
        <br>
        <br>
        <br>
        <br>
        <div><a href="/">About</a></div>
        <div><a href="https://x.com/145k4">@145k4</a></div>
		<div>hello@alaskahoffman.com</div>
    </nav>
</div>`

// Base template - main page layout with CSS styling
const baseTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>{{.Title}}</title>
    <style>
        body {
            font-family: Helvetica, Arial, sans-serif;
            font-size: 11px;
            background-color: white;
            margin: 0;
            padding: 0;
        }
        .container {
            display: flex;
            min-height: 100vh;
        }
        .sidebar {
            width: 150px;
            padding: 20px;
            background-color: transparent;
            border-right: 1px solid #ccc;
        }
          .main-content {
              flex: 1;
              padding: 80px 60px 20px 60px;
              max-width: 600px;
          }
          .mobile-menu-toggle {
              display: none;
          }
          @media (max-width: 768px) {
              .main-content {
                  padding: 40px 20px 20px 20px;
                  max-width: 100%;
              }
              .sidebar {
                  display: none;
                  position: fixed;
                  top: 0;
                  left: 0;
                  width: 200px;
                  height: 100vh;
                  background-color: white;
                  z-index: 1000;
                  padding: 20px;
                  border-right: 1px solid #ccc;
                  overflow-y: auto;
                  box-shadow: 2px 0 5px rgba(0,0,0,0.1);
                  font-size: 14px;
              }
              .sidebar p,
              .sidebar nav div,
              .sidebar nav a,
              .sidebar form input {
                  font-size: 14px;
              }
              .sidebar.mobile-open {
                  display: block;
              }
              .mobile-menu-toggle {
                  display: block;
                  padding: 15px;
                  border-bottom: 1px solid #ccc;
                  background-color: white;
                  position: sticky;
                  top: 0;
                  z-index: 999;
                  font-size: 14px;
              }
              .mobile-menu-toggle a {
                  font-weight: bold;
                  font-size: 14px;
              }
              .container {
                  flex-direction: column;
              }
          }
          .poem-list {
              line-height: 1.1;
              margin-top: 13px;
          }
          .poem-listing {
              margin-bottom: 2px;
          }
        input[type="text"] {
            width: 70px;
            height: 8px;
            padding: 3px;
            border: 1px solid black;
            border-radius: 0px;
            background-color: transparent;
            font-family: Helvetica, Arial, sans-serif;
            font-size: 11px;
        }
        form {
            margin-bottom: 13px;
        }
        a {
            color: black;
            text-decoration: none;
        }
        a:hover {
            color: black;
            text-decoration: underline;
        }
        .boma-content {
            font-size: 14px;
            line-height: 1.6;
            font-family: "Times New Roman", Times, serif;
        }
        .boma-content p {
            text-indent: 2em;
        }
    </style>
    <script>
        function toggleMobileMenu() {
            var sidebar = document.querySelector('.sidebar');
            sidebar.classList.toggle('mobile-open');
        }
    </script>
</head>
<body>
    <div class="mobile-menu-toggle">
        <a href="#" onclick="toggleMobileMenu(); return false;">Alaska Hoffman</a>
    </div>
    <div class="container">
        {{.Sidebar}}
        <div class="main-content">
            {{.Content}}
        </div>
    </div>
</body>
</html>`

// ============================================================================
// UTILITY FUNCTIONS
// ============================================================================

// formatBomaContent - formats story content with paragraph tags for proper indentation
func formatBomaContent(content string) string {
	// First, split on section headers to separate them from following text
	// Replace "\n\nI.\n" with a special marker, same for II. and III.
	content = strings.ReplaceAll(content, "\n\nI.\n", "\n\n__SECTION_I__\n")
	content = strings.ReplaceAll(content, "\n\nII.\n", "\n\n__SECTION_II__\n")
	content = strings.ReplaceAll(content, "\n\nIII\n", "\n\n__SECTION_III__\n")
	content = strings.ReplaceAll(content, "\n\nIII.\n", "\n\n__SECTION_III__\n")
	
	// Split by double newlines to get paragraphs
	paragraphs := strings.Split(content, "\n\n")
	var result strings.Builder
	
	for i, para := range paragraphs {
		para = strings.TrimSpace(para)
		if para == "" {
			continue
		}
		
		// Helper function to format content with indented line breaks
		formatWithIndentedBreaks := func(text string) string {
			lines := strings.Split(text, "\n")
			var formatted strings.Builder
			for i, line := range lines {
				if i > 0 {
					// Add indent for lines after the first
					formatted.WriteString("<br><span style=\"display: inline-block; text-indent: 2em;\">")
					formatted.WriteString(strings.TrimSpace(line))
					formatted.WriteString("</span>")
				} else {
					// First line - no indent needed (paragraph already has text-indent)
					if i < len(lines)-1 {
						formatted.WriteString(line)
					} else {
						formatted.WriteString(line)
					}
				}
			}
			return formatted.String()
		}
		
		// Check if this paragraph starts with a section marker
		if strings.HasPrefix(para, "__SECTION_I__") {
			result.WriteString("<p style=\"text-indent: 0; text-align: center;\">I.</p>\n")
			// Get the rest of the text after the marker
			remaining := strings.TrimPrefix(para, "__SECTION_I__")
			remaining = strings.TrimSpace(remaining)
			if remaining != "" {
				result.WriteString("<p>")
				result.WriteString(formatWithIndentedBreaks(remaining))
				result.WriteString("</p>")
			}
		} else if strings.HasPrefix(para, "__SECTION_II__") {
			result.WriteString("<p style=\"text-indent: 0; text-align: center;\">II.</p>\n")
			remaining := strings.TrimPrefix(para, "__SECTION_II__")
			remaining = strings.TrimSpace(remaining)
			if remaining != "" {
				result.WriteString("<p>")
				result.WriteString(formatWithIndentedBreaks(remaining))
				result.WriteString("</p>")
			}
		} else if strings.HasPrefix(para, "__SECTION_III__") {
			result.WriteString("<p style=\"text-indent: 0; text-align: center;\">III.</p>\n")
			remaining := strings.TrimPrefix(para, "__SECTION_III__")
			remaining = strings.TrimSpace(remaining)
			if remaining != "" {
				result.WriteString("<p>")
				result.WriteString(formatWithIndentedBreaks(remaining))
				result.WriteString("</p>")
			}
		} else {
			// Regular paragraph - format with indented breaks
			result.WriteString("<p>")
			result.WriteString(formatWithIndentedBreaks(para))
			result.WriteString("</p>")
		}
		
		// Add spacing between paragraphs (except for last one)
		if i < len(paragraphs)-1 {
			result.WriteString("\n")
		}
	}
	
	return result.String()
}

// searchPoems - searches through all poem JSON files for matching content
func searchPoems(query string) []struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Date     string `json:"date"`
	Category string `json:"category"`
	Location string `json:"location"`
	Content  string `json:"content"`
} {
	var results []struct {
		ID       int    `json:"id"`
		Title    string `json:"title"`
		Date     string `json:"date"`
		Category string `json:"category"`
		Location string `json:"location"`
		Content  string `json:"content"`
	}
	
	// Read all poem JSON files from embedded filesystem
	entries, err := staticFiles.ReadDir("public/poems")
	if err != nil {
		log.Printf("Error reading poem files: %v", err)
		return results
	}
	
	queryLower := strings.ToLower(query)
	
	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		data, err := staticFiles.ReadFile("public/poems/" + entry.Name())
		if err != nil {
			continue
		}
		
		var poem struct {
			ID       int    `json:"id"`
			Title    string `json:"title"`
			Date     string `json:"date"`
			Category string `json:"category"`
			Location string `json:"location"`
			Content  string `json:"content"`
		}
		
		if err := json.Unmarshal(data, &poem); err != nil {
			continue
		}
		
		// Search in title, content, category, and location
		if strings.Contains(strings.ToLower(poem.Title), queryLower) ||
		   strings.Contains(strings.ToLower(poem.Content), queryLower) ||
		   strings.Contains(strings.ToLower(poem.Category), queryLower) ||
		   strings.Contains(strings.ToLower(poem.Location), queryLower) {
			results = append(results, poem)
		}
	}
	return results
}

// ============================================================================
// HTTP HANDLERS
// ============================================================================

// Homepage handler - displays bio and navigation
func homeHandler(w http.ResponseWriter, r *http.Request) {
	sidebar, _ := template.New("sidebar").Parse(sidebarTemplate)
	base, _ := template.New("base").Parse(baseTemplate)
	
	var sidebarHTML strings.Builder
	sidebar.Execute(&sidebarHTML, struct{ Query string }{""})
	
	content := `
        <div class="bio">
            <p>Alaska Hoffman is a Michigander poet based in Brooklyn, New York.</p>
            <p>She is interested in themes of noise, futurism, permanence, hauntology, transition, repetition, and historicity.</p>
            <p>She has a B.A. in Creative Writing from Columbia University, and is a USMC veteran.</p>
            <p>This website serves primarily as an archive of her personal work.</p>
            <p>Alaska has also created under the names dovetail, ennen, and Archway Labs.</p>
            <br>
            <p><a href="https://x.com/145k4">@145k4</a></p>
            <p>hello@alaskahoffman.com</p>
        </div>`
	
	data := struct {
		Title   string
		Sidebar template.HTML
		Content template.HTML
	}{
		Title:   "Alaska Hoffman",
		Sidebar: template.HTML(sidebarHTML.String()),
		Content: template.HTML(content),
	}
	
	base.Execute(w, data)
}

// Search handler - searches through poem JSON files and displays results
func searchHandler(w http.ResponseWriter, r *http.Request) {
	query := strings.ToLower(r.URL.Query().Get("q"))
	
	// Search through poems
	var poemResults []struct {
		ID       int    `json:"id"`
		Title    string `json:"title"`
		Date     string `json:"date"`
		Category string `json:"category"`
		Location string `json:"location"`
		Content  string `json:"content"`
	}
	if query != "" {
		poemResults = searchPoems(query)
	}
	
	sidebar, _ := template.New("sidebar").Parse(sidebarTemplate)
	base, _ := template.New("base").Parse(baseTemplate)
	
	var sidebarHTML strings.Builder
	sidebar.Execute(&sidebarHTML, struct{ Query string }{query})
	
	content := fmt.Sprintf(`
        <h2>Search Results for "%s"</h2>
        <p>Found %d poems</p>
        
        %s
        
        <p><a href="/">← Back to Home</a></p>`, 
		query, len(poemResults),
		func() string {
			resultHTML := ""
			
			// Display poem results
			if len(poemResults) > 0 {
				for _, poem := range poemResults {
					// Truncate content for preview
					preview := poem.Content
					if len(preview) > 200 {
						preview = preview[:200] + "..."
					}
					resultHTML += fmt.Sprintf(`
            <div class="poem-result">
                <h4><a href="/poem/%d">%s</a></h4>
                <p><strong>Date:</strong> %s | <strong>Location:</strong> %s</p>
                <p class="poem-preview">%s</p>
            </div>`, poem.ID, poem.Title, poem.Date, poem.Location, preview)
				}
			} else {
				resultHTML = fmt.Sprintf(`<p>No poems found matching "%s"</p>`, query)
			}
			
			return resultHTML
		}())
	
	data := struct {
		Title   string
		Sidebar template.HTML
		Content template.HTML
	}{
		Title:   "Search Results - Alaska Hoffman",
		Sidebar: template.HTML(sidebarHTML.String()),
		Content: template.HTML(content),
	}
	
	base.Execute(w, data)
}

// Poem handler - displays individual poem pages from JSON files
func poemHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/poem/")
	poemID := strings.TrimSuffix(path, ".json")
	
	// Read the JSON file from embedded filesystem
	filePath := fmt.Sprintf("public/poems/poem-%s.json", poemID)
	data, err := staticFiles.ReadFile(filePath)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	
	// Parse the JSON
	var poem struct {
		ID       int    `json:"id"`
		Title    string `json:"title"`
		Date     string `json:"date"`
		Category string `json:"category"`
		Location string `json:"location"`
		Content  string `json:"content"`
	}
	
	if err := json.Unmarshal(data, &poem); err != nil {
		http.Error(w, "Error parsing poem data", http.StatusInternalServerError)
		return
	}
	
	sidebar, _ := template.New("sidebar").Parse(sidebarTemplate)
	base, _ := template.New("base").Parse(baseTemplate)
	
	var sidebarHTML strings.Builder
	sidebar.Execute(&sidebarHTML, struct{ Query string }{""})
	
	content := fmt.Sprintf(`
        <h4>%s</h4>
        
        <div class="poem-content">
            %s
        </div>
        <br>
        <br>
        <p>%s // %s</p>`,
		poem.Title, strings.ReplaceAll(poem.Content, "\n", "<br>"), 
		poem.Location, poem.Date)
	
	pageData := struct {
		Title   string
		Sidebar template.HTML
		Content template.HTML
	}{
		Title:   fmt.Sprintf("%s - Alaska Hoffman", poem.Title),
		Sidebar: template.HTML(sidebarHTML.String()),
		Content: template.HTML(content),
	}
	
	base.Execute(w, pageData)
}

// Boma handler - displays the Boma short story
func bomaHandler(w http.ResponseWriter, r *http.Request) {
	// Read the JSON file from embedded filesystem
	data, err := staticFiles.ReadFile("public/boma.json")
	if err != nil {
		http.NotFound(w, r)
		return
	}
	
	// Parse the JSON
	var story struct {
		Title    string `json:"title"`
		Date     string `json:"date"`
		Location string `json:"location"`
		Content  string `json:"content"`
	}
	
	if err := json.Unmarshal(data, &story); err != nil {
		http.Error(w, "Error parsing story data", http.StatusInternalServerError)
		return
	}
	
	sidebar, _ := template.New("sidebar").Parse(sidebarTemplate)
	base, _ := template.New("base").Parse(baseTemplate)
	
	var sidebarHTML strings.Builder
	sidebar.Execute(&sidebarHTML, struct{ Query string }{""})
	
	// Format content with paragraph tags for proper indentation
	formattedContent := formatBomaContent(story.Content)
	
	content := fmt.Sprintf(`
        <h4>%s</h4>
        
        <div class="poem-content boma-content">
            %s
        </div>
        <br>
        <br>
        <p>%s // %s</p>`,
		story.Title, formattedContent, 
		story.Location, story.Date)
	
	pageData := struct {
		Title   string
		Sidebar template.HTML
		Content template.HTML
	}{
		Title:   fmt.Sprintf("%s - Alaska Hoffman", story.Title),
		Sidebar: template.HTML(sidebarHTML.String()),
		Content: template.HTML(content),
	}
	
	base.Execute(w, pageData)
}

// Poetry handler - displays listing of all poems
func poetryHandler(w http.ResponseWriter, r *http.Request) {
	// Read all poem JSON files from embedded filesystem
	entries, err := staticFiles.ReadDir("public/poems")
	if err != nil {
		log.Printf("Error reading poem files: %v", err)
		http.Error(w, "Error reading poems", http.StatusInternalServerError)
		return
	}
	
	var poems []struct {
		ID       int    `json:"id"`
		Title    string `json:"title"`
		Date     string `json:"date"`
		Category string `json:"category"`
		Location string `json:"location"`
		Content  string `json:"content"`
	}
	
	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		data, err := staticFiles.ReadFile("public/poems/" + entry.Name())
		if err != nil {
			continue
		}
		
		var poem struct {
			ID       int    `json:"id"`
			Title    string `json:"title"`
			Date     string `json:"date"`
			Category string `json:"category"`
			Location string `json:"location"`
			Content  string `json:"content"`
		}
		
		if err := json.Unmarshal(data, &poem); err != nil {
			continue
		}
		
		poems = append(poems, poem)
	}
	
	sidebar, _ := template.New("sidebar").Parse(sidebarTemplate)
	base, _ := template.New("base").Parse(baseTemplate)
	
	var sidebarHTML strings.Builder
	sidebar.Execute(&sidebarHTML, struct{ Query string }{""})
	
	content := fmt.Sprintf(`
        <div class="poem-list">
        %s
        </div>
        
        <p><a href="/">← Back to Home</a></p>`, 
		func() string {
			if len(poems) > 0 {
				resultHTML := ""
				for _, poem := range poems {
					resultHTML += fmt.Sprintf(`
            <div class="poem-listing">
                <a href="/poem/%d">%s</a>
            </div>`, poem.ID, poem.Title)
				}
				return resultHTML
			}
			return `<p>No poems found in the archive.</p>`
		}())
	
	pageData := struct {
		Title   string
		Sidebar template.HTML
		Content template.HTML
	}{
		Title:   "All Poems - Alaska Hoffman",
		Sidebar: template.HTML(sidebarHTML.String()),
		Content: template.HTML(content),
	}
	
	base.Execute(w, pageData)
}

// ============================================================================
// MAIN HANDLER FUNCTION
// ============================================================================

// Handler is the main entry point for Vercel Go functions
func Handler(w http.ResponseWriter, r *http.Request) {
	// Initialize random seed
	rand.Seed(time.Now().UnixNano())
	
	// Handle static files first
	if strings.HasPrefix(r.URL.Path, "/static/") {
		filePath := strings.TrimPrefix(r.URL.Path, "/static/")
		data, err := staticFiles.ReadFile("public/" + filePath)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		
		// Set appropriate content type and caching headers
		if strings.HasSuffix(filePath, ".webp") {
			w.Header().Set("Content-Type", "image/webp")
			w.Header().Set("Cache-Control", "public, max-age=31536000")
		} else if strings.HasSuffix(filePath, ".json") {
			w.Header().Set("Content-Type", "application/json")
		}
		
		w.Write(data)
		return
	}
	
	// Route handling
	switch {
	case r.URL.Path == "/" || r.URL.Path == "":
		homeHandler(w, r)
	case r.URL.Path == "/search":
		searchHandler(w, r)
	case strings.HasPrefix(r.URL.Path, "/poem/"):
		poemHandler(w, r)
	case r.URL.Path == "/poetry":
		poetryHandler(w, r)
	case r.URL.Path == "/boma":
		bomaHandler(w, r)
	default:
		http.NotFound(w, r)
	}
}

