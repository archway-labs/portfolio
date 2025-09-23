package handler

// ============================================================================
// IMPORTS
// ============================================================================
import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

// ============================================================================
// DATA STRUCTURES
// ============================================================================
// No data structures needed - all data is served from JSON files

// ============================================================================
// MAIN SERVER SETUP
// ============================================================================

// Global mux for routing
var mux *http.ServeMux

// Initialize the server routes
func init() {
	rand.Seed(time.Now().UnixNano())
	
	mux = http.NewServeMux()
	
	// Route handlers
	mux.HandleFunc("/", homeHandler)           // Homepage with bio
	mux.HandleFunc("/search", searchHandler)   // Search poems functionality
	mux.HandleFunc("/poem/", poemHandler)      // Individual poem pages
	mux.HandleFunc("/poetry", poetryHandler)   // All poems listing page
	
	// Static files are served automatically by Vercel from public/ directory
	mux.HandleFunc("/debug", debugHandler) // Debug endpoint to check file access
}

// Handler function for Vercel
func Handler(w http.ResponseWriter, r *http.Request) {
	mux.ServeHTTP(w, r)
}

// Main function for local development
func main() {
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
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
        <s><div><a href="/search?q=">Calligraphy</a></div></s>
        <s><div><a href="/search?q=">Paintings</a></div></s>
        <s><div><a href="/search?q=">Photography</a></div></s>
        <s><div><a href="/search?q=">Media Art</a></div></s>
        <br>
        <s><div><a href="/search?q=">Capsule 21 (2022)</a></div></s>
        <s><div><a href="/search?q=">C21 Babylon (2022)</a></div></s>
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
        <div><a href="https://x.com/145k4">@145k4</a></div>
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
            background-image: url('/archbgs-01.webp');
            background-size: cover;
            background-position: center;
            background-repeat: no-repeat;
            background-attachment: fixed;
            margin: 0;
            padding: 0;
        }
        body::before {
            content: '';
            position: fixed;
            top: 0;
            left: 0;
            width: 100%;
            height: 100%;
            background-color: rgba(255, 255, 255, 0.97);
            z-index: -1;
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
              padding: 80px 20px 20px 60px;
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
    </style>
</head>
<body>
    <div class="container">
        {{.Sidebar}}
        <div class="main-content">
            {{.Content}}
        </div>
    </div>
</body>
</html>`

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

// ============================================================================
// UTILITY FUNCTIONS
// ============================================================================

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
	
	// Get the base URL from environment or use localhost for development
	baseURL := os.Getenv("VERCEL_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	} else {
		baseURL = "https://" + baseURL
	}
	
	queryLower := strings.ToLower(query)
	
	// Try to fetch each poem file via HTTP
	for i := 1; i <= 50; i++ { // Assuming we have up to 50 poems
		url := fmt.Sprintf("%s/poems/poem-%d.json", baseURL, i)
		resp, err := http.Get(url)
		if err != nil || resp.StatusCode != 200 {
			continue
		}
		
		data, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
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
	
	// Get the base URL from environment or use localhost for development
	baseURL := os.Getenv("VERCEL_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	} else {
		baseURL = "https://" + baseURL
	}
	
	// Fetch the poem JSON file via HTTP
	url := fmt.Sprintf("%s/poems/poem-%s.json", baseURL, poemID)
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != 200 {
		http.NotFound(w, r)
		return
	}
	defer resp.Body.Close()
	
	data, err := ioutil.ReadAll(resp.Body)
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

// Poetry handler - displays listing of all poems
func poetryHandler(w http.ResponseWriter, r *http.Request) {
	// Get the base URL from environment or use localhost for development
	baseURL := os.Getenv("VERCEL_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	} else {
		baseURL = "https://" + baseURL
	}
	
	var poems []struct {
		ID       int    `json:"id"`
		Title    string `json:"title"`
		Date     string `json:"date"`
		Category string `json:"category"`
		Location string `json:"location"`
		Content  string `json:"content"`
	}
	
	// Try to fetch each poem file via HTTP
	for i := 1; i <= 50; i++ { // Assuming we have up to 50 poems
		url := fmt.Sprintf("%s/poems/poem-%d.json", baseURL, i)
		resp, err := http.Get(url)
		if err != nil || resp.StatusCode != 200 {
			continue
		}
		
		data, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
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

// Debug handler to check file access in Vercel environment
func debugHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	
	// Check current working directory
	wd, _ := os.Getwd()
	fmt.Fprintf(w, "Working Directory: %s\n", wd)
	
	// Check if files exist
	files := []string{
		"public/poems/poem-1.json",
		"public/archbgs-01.webp",
		"archbgs-01.webp",
		"poems/poem-1.json",
	}
	
	for _, file := range files {
		if _, err := os.Stat(file); err == nil {
			fmt.Fprintf(w, "✓ %s exists\n", file)
		} else {
			fmt.Fprintf(w, "✗ %s not found: %v\n", file, err)
		}
	}
	
	// List directory contents
	fmt.Fprintf(w, "\nDirectory listing:\n")
	if dirs, err := os.ReadDir("."); err == nil {
		for _, dir := range dirs {
			fmt.Fprintf(w, "- %s\n", dir.Name())
		}
	}
}

