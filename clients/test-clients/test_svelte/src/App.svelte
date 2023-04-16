<script context="module" lang="ts">
  import { writable } from "svelte/store";

  import {
    signInWithRedirect,
    GoogleAuthProvider,
    signOut,
    getAuth,
    getRedirectResult,
  } from "firebase/auth";
  import type { User } from "firebase/auth";
  import { initializeApp } from "firebase/app";
  import type { PostOverviewCollection, Post } from "./lib/DataInterface";
  import { LoadPosts, LoadPost } from "./lib/DataService";
  import type { AppUser } from "./lib/DataInterface";

  const googleProvider = new GoogleAuthProvider();

  // Your web app's Firebase configuration
  const firebaseConfig = {
    apiKey: "AIzaSyDWh5sExqNSMsT8Jj6-0q01j6KWL_UmX48",
    authDomain: "cloud32x.firebaseapp.com",
    projectId: "cloud32x",
    storageBucket: "cloud32x.appspot.com",
    messagingSenderId: "323709580283",
    appId: "1:323709580283:web:860fcce17eadee754d915e",
  };

  // Initialize Firebase & Firebase Auth
  const app = initializeApp(firebaseConfig);
  const auth = getAuth(app);

  let currentUser: User = undefined;

  export const user = writable(undefined);
  let userSignedIn: boolean = primitiveToBoolean(
    localStorage.getItem("UserSignedIn")
  );

  export const url = writable("/");
  if (userSignedIn) navigate("/home");

  // Navigate and route in the app
  export function navigate(path: string) {
    console.log("navigate " + path);

    history.pushState(undefined, "", path);

    url.update((n) => path);
  }

  // User functions
  export function getIdToken(): Promise<string> {
    return new Promise<string>((resolve, reject) => {
      if (currentUser) {
        currentUser
          .getIdToken(/* forceRefresh */ true)
          .then(function (idToken) {
            resolve(idToken);
          });
      } else {
        reject("No current user.");
      }
    });
  }

  export function signInWithGoogle() {
    const auth = getAuth();
    signInWithRedirect(auth, googleProvider);
  }

  export function signUserOut() {
    signOut(auth)
      .then(() => {
        // Sign-out successful.
        currentUser = undefined;
        user.update((n) => undefined);

        navigate("/");
      })
      .catch((error) => {
        // An error happened.
      });
  }

  auth.onAuthStateChanged((u: User) => {
    // if u is undefined, means we don't know user state
    // if u is null, means user is signed out
    // if u is an object, means user is signed in
    currentUser = u;
    // console.log("Setting user: " + currentUser.email);
    // let newUser: AppUser = {
    //   displayName: currentUser.displayName,
    //   email: currentUser.email,
    //   photoURL: currentUser.photoURL,
    //   phoneNumber: currentUser.phoneNumber,
    //   providerId: currentUser.providerId,
    //   uid: currentUser.uid,
    //   emailVerified: currentUser.emailVerified,
    //   isAnonymous: currentUser.isAnonymous,
    //   refreshToken: currentUser.refreshToken,
    // };

    if (!u) {
      userSignedIn = false;
      localStorage.setItem("UserSignedIn", "false");
    } else {
      userSignedIn = true;
      localStorage.setItem("UserSignedIn", "true");
    }

    user.update((n) => currentUser as AppUser);

    if (currentUser) navigate("/home");
    else navigate("/");
  });

  // This is only useful for diagnosing Firebase Auth problems, like when 3rd party cookies are disabled
  // getRedirectResult(auth)
  //   .then((result) => {
  //     console.log("getRedirectResult success")
  //     if (result.user)
  //       console.log("getRedirectResult_userset")
  //     else
  //       console.log("getRedirectResult_nouserset")

  //   }).catch((error) => {
  //     // Handle Errors here.
  //     console.log("getRedirectResult error: " + JSON.stringify(error))
  //   });

  console.log("setting popstate event");
  addEventListener("popstate", (event) => {
    console.log("popstate event:");
    console.log(location.href);

    const url = new URL(location.href);
    console.log(url.pathname);
    navigate(url.pathname);
  });

  function primitiveToBoolean(
    value: string | number | boolean | null | undefined
  ): boolean {
    if (typeof value === "string") {
      return value.toLowerCase() === "true" || !!+value; // here we parse to number first
    }

    return !!value;
  }
</script>

<script lang="ts">
  // import { Router, Route, navigate } from "svelte-routing";
  import { get } from "svelte/store";

  import Header from "./lib/Header.svelte";
  import Welcome from "./lib/Welcome.svelte";
  import Posts from "./lib/Posts.view.svelte";
  import NewPost from "./lib/Post.new.svelte";
  import ViewPost from "./lib/Post.view.svelte";
  import EditPost from "./lib/Post.edit.svelte";

  let localUrl: string = "";
  let rootPath: string = "";

  let pathParts: string[] = [];

  let allPosts: PostOverviewCollection = {};
  let currentPostId: string = "";
  let currentPost: Post;

  let localUser: AppUser = undefined;

  user.subscribe((value: AppUser) => {
    localUser = value;
  });

  url.subscribe((value) => {
    console.log("New url: " + value);
    localUrl = value;

    const parts: string[] = localUrl.split("/");
    parts.shift();

    switch (parts[0]) {
      case "home":
        rootPath = parts[0];
        LoadPosts().then((posts) => {
          allPosts = posts;
        });
        break;
      case "posts":
        currentPost = undefined;
        console.log("loading post " + parts[1]);
        rootPath = parts[0];
        LoadPost(parts[1]).then((post) => {
          console.log("finished loading post");
          currentPost = post;
          // parts.shift();
          // pathParts = parts;
        });
        break;
      default:
        rootPath = parts[0];
    }
  });
</script>

<main>
  <div class="ac">
    {#if rootPath == ""}
      <Header />
      <Welcome />
    {:else if rootPath == "home"}
      <Header />
      <Posts posts={allPosts} />
    {:else if rootPath == "posts"}
      <ViewPost id={pathParts[0]} post={currentPost} />
    {:else if rootPath == "new"}
      <Header small={true} />
      <NewPost />
    {:else}
      <h1>404</h1>
      <div>{rootPath}</div>
    {/if}
  </div>
</main>

<style>
  .ac {
    margin: 0px;
  }
</style>
