# post
Post resource for REST API that provides concurrency safe accessor and mutator functions using a SQLite3 DB for persistence.

## Quick Start
See https://github.com/jecolon/restsrv for a sample REST API server that uses this package.

```go
// Init must be called before using the rest of the package API.
if err := post.Init(); err != nil {
  log.Fatalf("post.Init() error: %v", err)
}
  
// A Post instance
p1 := post.Post{
  UserId: 1,
  Title: "The title",
  Body: "The body of the post.",
}

// Save the new post
post.New(p1)

// Get a post. id is an int. found is false if post not found.
p2, found := post.Get(id)

// List all posts. posts is a []Post.
posts := post.List()

// Update an existing post.
p1.Body = "A new body!"
post.Put(p1)

// Delete a post.
post.Del(p1.Id)
```
