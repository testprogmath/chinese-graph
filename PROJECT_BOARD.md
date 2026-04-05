# 📋 Project Board Setup

## GitHub Projects Configuration

### 1. Create Project Board
Go to: https://github.com/users/YOUR_USERNAME/projects/new

**Settings:**
- Name: Chinese Learning Graph
- Description: Visual graph-based Chinese learning application
- Visibility: Private (or Public)
- Template: Board

### 2. Columns Setup
```
📋 Backlog → 🎯 Ready → 🚧 In Progress → 👀 Review → ✅ Done
```

### 3. Automation Rules

**Backlog**
- New issues automatically added here
- Sorted by priority label

**Ready**
- Move here when:
  - Requirements clear
  - Dependencies resolved
  - Assigned to milestone

**In Progress**
- Auto-move when:
  - Issue assigned
  - PR created with `fixes #issue`
  - Label `status:in-progress` added

**Review**
- Auto-move when:
  - PR ready for review
  - Label `status:review` added

**Done**
- Auto-move when:
  - Issue closed
  - PR merged

### 4. Quick Commands in Comments

```bash
/assign @username     # Assign issue
/priority P0          # Set priority (P0-P3)
/estimate XL          # Set size estimate
/status blocked       # Update status
/milestone MVP        # Add to milestone
```

### 5. Current Backlog

#### 🔴 High Priority (MVP)
- [ ] Frontend setup with React + TypeScript
- [ ] Graph visualization with React Flow
- [ ] User authentication
- [ ] Word search functionality
- [ ] Basic learning stats

#### 🟡 Medium Priority (v1.0)
- [ ] Mobile app (React Native)
- [ ] Spaced repetition algorithm
- [ ] Export to Anki
- [ ] Batch word import
- [ ] User preferences

#### 🟢 Low Priority (Future)
- [ ] Gamification system
- [ ] Social features
- [ ] AI-powered suggestions
- [ ] Voice recognition
- [ ] Handwriting practice

### 6. Milestones

**MVP (Target: 2 weeks)**
- Core graph visualization
- Basic CRUD operations
- Search functionality
- Docker deployment

**Mobile v1 (Target: 4 weeks)**
- React Native app
- Offline support
- Push notifications

**Gamification (Target: 6 weeks)**
- Points system
- Achievements
- Daily streaks
- Leaderboards

### 7. Weekly Sync Format

```markdown
## Week of [Date]

### ✅ Completed
- Task 1
- Task 2

### 🚧 In Progress
- Task 3 (70% done)
- Task 4 (just started)

### 🚫 Blocked
- Task 5 (waiting for X)

### 📅 Next Week
- Priority 1
- Priority 2

### 💭 Notes
- Any important decisions
- Technical debt to address
```

## How I'll Use This

When you request a task, I will:

1. **Create GitHub Issue** with:
   - Clear description
   - Acceptance criteria
   - Size estimate
   - Priority label
   - Component labels

2. **Update Status** as I work:
   - Move to "In Progress" when starting
   - Add comments with progress
   - Link PRs to issues
   - Move to "Done" when complete

3. **Track Dependencies** by:
   - Linking related issues
   - Using `blocked by #123` in comments
   - Updating milestone progress

## Example Issue Creation

```bash
gh issue create \
  --title "Implement graph visualization with React Flow" \
  --body "### Description
Visual graph component showing word relationships

### Acceptance Criteria
- [ ] Interactive node dragging
- [ ] Zoom/pan controls
- [ ] Different edge types for relationship types
- [ ] Node clustering for large graphs
- [ ] Search highlighting

### Technical Notes
- Use React Flow library
- Implement virtual scrolling for performance
- Cache rendered graphs in Redis" \
  --label "frontend,type:feature,priority:high,size:L" \
  --milestone "MVP" \
  --project "Chinese Learning App"
```

## Access Your Board

1. Go to your GitHub profile
2. Click "Projects" tab
3. Open "Chinese Learning Graph"
4. Or direct link: `https://github.com/users/[YOUR_USERNAME]/projects/1`

Ready to start tracking! Should I create the first batch of issues for the MVP?