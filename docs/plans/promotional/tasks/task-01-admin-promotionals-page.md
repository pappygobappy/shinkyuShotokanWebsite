# Task 1: Admin Promotionals Page — Route + Handler + Skeleton Template

**Files:**
- Modify: `routes/admin.go`
- Create: `handlers/promotional.go`
- Create: `templates/admin_promotionals.html`
- Modify: `templates/adminPage.html`

**Interfaces:**
- Consumes: none (first task)
- Produces: `AdminPromotionalsPage` handler function; `/admin/promotionals` route; `admin_promotionals` template define block

**Steps:**

- [ ] **Step 1: Add the route to admin routes**

In `routes/admin.go`, add this block after the existing Gear management routes (after line 92, before the closing `}`):

```go
	// Promotional management
	adminRoutes.Get("/promotionals", handlers.AdminPromotionalsPage)
```

- [ ] **Step 2: Create the handler file**

Create `handlers/promotional.go` with this exact content:

```go
package handlers

import (
	"shinkyuShotokan/structs"
	"shinkyuShotokan/utils"
)

func AdminPromotionalsPage(c *fiber.Ctx) error {
	page := fiber.Map{
		"Page":  structs.Page{PageName: "Promotionals"},
		"Tabs":  utils.CurrentTabs(),
	}
	return c.Render("admin_promotionals", page)
}
```

- [ ] **Step 3: Create the skeleton template**

Create `templates/admin_promotionals.html` with this exact content:

```html
{{ define "admin_promotionals" }}
<div class="text-3xl border-b py-5">
    <div class="grow">
        <div>Promotionals</div>
    </div>
</div>
<div class="overflow-x-auto">
    <table class="table table-sm table-zebra" style="min-width: 500px;">
        <thead>
            <tr>
                <th>Name</th>
                <th>Date</th>
                <th>SubType</th>
                <th>Location</th>
                <th>Applicants</th>
                <th>Actions</th>
            </tr>
        </thead>
        <tbody>
            {{ range .Promotionals }}
            <tr>
                <td>{{ .Name }}</td>
                <td>{{ .Date.Format "2006-01-02" }}</td>
                <td>{{ .SubType }}</td>
                <td>{{ .Location }}</td>
                <td>0</td>
                <td>
                    <div class="flex flex-row">
                        <button class="btn btn-sm btn-primary">Edit</button>
                        <button class="btn btn-sm btn-error">Delete</button>
                    </div>
                </td>
            </tr>
            {{ end }}
        </tbody>
    </table>
</div>
{{ end }}
```

- [ ] **Step 4: Wire the template into adminPage.html**

In `templates/adminPage.html`, find the template conditional block (around line 11-21) that looks like:

```html
{{ if eq .Page.PageName "Locations"}}
    {{ template "locations" . }}
{{ else if eq .Page.PageName "Event Templates"}}
    {{ template "admin_event_template" . }} 
{{ else if eq .Page.PageName "Users"}}
    {{ template "users" . }}
{{ else if eq .Page.PageName "User Profile"}}
    {{ template "userProfile" . }}
{{ else if eq .Page.PageName "Carousel Images"}}
    {{ template "admin_carousel_images" . }}
{{ end }}
```

Add a new `else if` branch before the final `{{ end }}`:

```html
{{ else if eq .Page.PageName "Promotionals"}}
    {{ template "admin_promotionals" . }}
```

- [ ] **Step 5: Add sidebar link in adminPage.html**

In `templates/adminPage.html`, find the sidebar menu (around line 25-58). Add a new `<li>` block after the "Event Templates" link (after line 44):

```html
<li>
    <a href="/admin/promotionals">
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" width="24" height="24" stroke="currentColor" stroke-width="2"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="16" y1="13" x2="8" y2="13"/><line x1="16" y1="17" x2="8" y2="17"/><polyline points="10 9 9 9 8 9"/></svg>
        Promotionals
    </a>
</li>
```

- [ ] **Step 6: Add navbar dropdown link**

In `templates/adminPage.html`, find the navbar dropdown menu section (the `<ul class="dropdown-menu ...">` or similar nav links). Add a link to `/admin/promotionals` alongside existing admin links. Look for the navigation area that contains links like `/admin/locations`, `/admin/event-templates`, etc. and add:

```html
<a href="/admin/promotionals">Promotionals</a>
```

- [ ] **Step 7: Verify with go build**

Run: `go build ./...`
Expected: PASS (no errors)

- [ ] **Step 8: Manual verification**

Run: `go run main.go`
Open browser to `http://localhost:8080/admin/promotionals` (or whatever port the app uses)
Expected: Page renders with "Promotionals" heading, empty table with correct column headers, sidebar shows "Promotionals" link, navbar dropdown includes "Promotionals"
