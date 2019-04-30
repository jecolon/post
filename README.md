# post
Post resource for REST API that provides concurrency safe accessor and mutator functions.

## Quick Start
See https://github.com/jecolon/restsrv for a sample REST API server that uses this package.

```go
// Start the monitor goroutine
post.Start()

// A Post instance
p1 := post.Post{
  Id: post.NewId(),
  UserId: 1,
  Title: "The title",
  Body: "The body of the post.",
}

// Get a post. id is an int.
p2 := post.Get(id)

// List all posts. posts is a []Post.
posts := post.List()

// Save a new post.
post.Add(p1)

// Update an existing post.
p1.Body = "A new body!"
post.Set(p1)

// Delete a post.
post.Del(p1.Id)

// Stop the monitor goroutine.
post.Stop()
```
