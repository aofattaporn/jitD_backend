package controllers

type Post struct {
	ID      int
	Content string
	UserID  int
}

func ReccomendPost() {
	// // Generate some example posts
	// var posts []*Post
	// for i := 1; i <= 10; i++ {
	// 	post := &Post{ID: i, Content: fmt.Sprintf("Post %d content", i), UserID: i % 3}
	// 	posts = append(posts, post)
	// }

	// // // Build a k-d tree from the feature vectors of the posts
	// var tree *kdtree.KDTree

	// for i, _ := range posts {
	// 	fmt.Print(i)
	// 	// tree.Insert(points.NewPoint([]float64{12, 4, 6}, i))
	// }

	// // // Find the k-nearest neighbors of a new post
	// pointNewPost := points.NewPoint([]float64{12, 4, 6}, 0)

	// k := 3
	// neighbors := tree.KNN(pointNewPost, k)

	// // // Recommend the posts based on the k-nearest neighbors
	// var recommendedPosts []*Post
	// fmt.Println("==================================")
	// for _, neighbor := range neighbors {
	// 	post := neighbor.Dimensions()
	// 	fmt.Println(post)
	// 	// recommendedPosts = append(recommendedPosts, post)
	// }
	// fmt.Println("==================================")

	// fmt.Printf("Recommended posts: %v\n", recommendedPosts)
}
