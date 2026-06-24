# Task 8: Wire Promotional Link into Admin Sidebar

**Files:**
- Modify: `templates/adminPage.html`

**Interfaces:**
- Consumes: none (Task 1 already added sidebar + navbar links — this task is a no-op if Task 1 was completed correctly)
- Produces: no-op

**Note:** Task 1 already adds the sidebar link and navbar dropdown link. If those steps were completed, skip this task entirely. If the sidebar/navbar links were not added in Task 1, follow Step 1 below.

**Steps:**

- [ ] **Step 1 (only if Task 1 sidebar/navbar steps were skipped): Add sidebar link**

In `templates/adminPage.html`, add to the sidebar menu:

```html
<li>
    <a href="/admin/promotionals">
        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" width="24" height="24" stroke="currentColor" stroke-width="2"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="16" y1="13" x2="8" y2="13"/><line x1="16" y1="17" x2="8" y2="17"/><polyline points="10 9 9 9 8 9"/></svg>
        Promotionals
    </a>
</li>
```

- [ ] **Step 2 (only if Task 1 navbar step was skipped): Add navbar dropdown link**

In the navbar dropdown menu in `templates/adminPage.html`, add:
```html
<a href="/admin/promotionals">Promotionals</a>
```
