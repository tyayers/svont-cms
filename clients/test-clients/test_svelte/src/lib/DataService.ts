import { getIdToken } from "../App.svelte"
import type { PostOverviewCollection, Post, SearchResult } from "./DataInterface";
import logo from "../assets/bayer_gcloud_logo.png";

export let LogoPath: string = logo;

export async function SearchPosts(input: string): Promise<SearchResult[]> {
  return new Promise<SearchResult[]>((resolve, reject) => {
    if (input) {
      getIdToken().then(function (idToken) {
        fetch(
          import.meta.env.VITE_CMS_SERVICE +
            "/posts/search?q=" +
            input,
          {
            method: "GET",
            headers: {
              Accept: "application/json",
              Authorization: "Bearer " + idToken,
            },
          }
        )
          .then((response) => {
            return response.json();
          })
          .then((data: SearchResult[]) => {
            resolve(data);
          });
      });
    } else resolve([]);
  });
}

export async function LoadPosts(): Promise<PostOverviewCollection> {
  return new Promise<PostOverviewCollection>((resolve, reject) => {
    getIdToken().then((idToken) => {
      fetch(import.meta.env.VITE_CMS_SERVICE + "/posts", {
          method: "get",
          headers: {
            Accept: "application/json",
            Authorization: "Bearer " + idToken,
          },
        })
          .then((response) => {
            return response.json();
          })
          .then((data: PostOverviewCollection) => {
            resolve(data);
          });
        });
    });
}

export async function LoadPost(postId: string): Promise<Post> {
  return new Promise<Post>((resolve, reject) => {
    getIdToken().then((idToken) => {
      fetch(import.meta.env.VITE_CMS_SERVICE + "/posts/" + postId, {
        method: "get",
        headers: {
          Accept: "application/json",
          Authorization: "Bearer " + idToken,
        },
      })
        .then((response) => {
          return response.json();
        })
        .then((data: Post) => {
          resolve(data);
        });
      });
  });
}