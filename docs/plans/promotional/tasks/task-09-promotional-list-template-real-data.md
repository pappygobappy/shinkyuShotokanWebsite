# Task 9: Promotional List Template (real data)

**Files:**
- Modify: `templates/admin_promotionals.html`

**Interfaces:**
- Consumes: `queries.GetCheckedInCountByPromotionalID()` from Task 4
- Produces: fully functional list page with real data, edit/delete links, applicant counts

**Steps:**

- [ ] **Step 1: Update the handler to pass applicant counts**

In `handlers/promotional.go`, update `AdminPromotionalsPage` to compute checked-in counts:

```go
func AdminPromotionalsPage(c *fiber.Ctx) error {
	promotionals := queries.GetAllPromotionals()

	// Compute checked-in counts for each promotional
	counts := make(map[uint]int)
	for _, p := range promotionals {
		count, _ := queries.GetCheckedInCountByPromotionalID(p.ID)
		counts[p.ID] = count
	}

	var user *models.User
	if u := c.Locals("user"); u != nil {
		u := u.(models.User)
		if u.Type != "" {
			user = &u
		}
	}
	page := fiber.Map{
		"Page":           structs.Page{PageName: "Promotionals"},
		"Tabs":           utils.CurrentTabs(),
		"Promotionals":   promotionals,
		"CheckedInCount": counts,
		"user":           user,
	}
	return c.Render("admin_promotionals", page)
}
```

- [ ] **Step 2: Update the template to show real data**

Replace the entire content of `templates/admin_promotionals.html` with:

```html
{{ define "admin_promotionals" }}
<div class="text-3xl border-b py-5 flex flex-col md:flex-row md:items-center">
    <div class="grow">
        <div>Promotionals</div>
    </div>
    <div>
        <a href="/admin/promotionals/new" class="btn btn-primary btn-sm">Add Promotional</a>
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
                <th>Checked In</th>
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
                <td>
                    <a href="/admin/promotionals/{{ .ID }}/applicants" class="link link-primary">
                        0
                    </a>
                </td>
                <td>{{ index $.CheckedInCount .ID }}/0</td>
                <td>
                    <div class="flex flex-row gap-1">
                        <a href="/admin/promotionals/{{ .ID }}/edit" class="btn btn-sm btn-primary">Edit</a>
                        <button hx-delete="/admin/promotionals/{{ .ID }}" hx-confirm="Are you sure?" class="btn btn-sm btn-error">Delete</button>
                    </div>
                </td>
            </tr>
            {{ end }}
        </tbody>
    </table>
</div>
{{ end }}
```

Note: The applicant count column shows `0` as a placeholder. Task 10 (applicant handlers) will wire up the count properly.

- [ ] **Step 3: Verify with go build**

Run: `go build ./...`
Expected: PASS

- [ ] **Step 4: Manual verification**

Run: `go run main.go`
Navigate to `/admin/promotionals`
Expected: Table shows real promotional data (if any exist in DB); "Add Promotional" button visible; edit/delete links functional
