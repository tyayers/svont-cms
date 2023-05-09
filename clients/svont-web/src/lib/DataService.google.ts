import { goto, invalidate } from "$app/navigation";
import { browser } from "$app/environment";

import {
  signInWithRedirect,
  GoogleAuthProvider,
  signOut,
  getAuth,
  getRedirectResult,
} from "firebase/auth";
import type { User } from "firebase/auth";
import { initializeApp } from "firebase/app";

import {
  type SearchResult,
  type PostOverview,
  type PostOverviewCollection,
  type Post,
  type DataService,
  UserState,
  SignInProvider,
  EventType,
  type AppEvent,
  type PostComment,
  type Metadata,
} from "./DataInterface";
import { writable, type StartStopNotifier, type Writable } from "svelte/store";
import type { AppUser } from "./DataInterface";

// This class implements the DataService interface with pure test data,
// letting the complete application run in a local browser session wo
// a backend.
export class DataServiceGoogle implements DataService {
  googleProvider = new GoogleAuthProvider();

  // Your web app's Firebase configuration
  firebaseConfig = {
    apiKey: "AIzaSyDWh5sExqNSMsT8Jj6-0q01j6KWL_UmX48",
    authDomain: "cloud32x.firebaseapp.com",
    projectId: "cloud32x",
    storageBucket: "cloud32x.appspot.com",
    messagingSenderId: "323709580283",
    appId: "1:323709580283:web:860fcce17eadee754d915e",
  };

  // Initialize Firebase & Firebase Auth
  app = initializeApp(this.firebaseConfig);
  auth = getAuth(this.app);
  currentUser: User = undefined;

  localUser: AppUser = undefined;
  localUserState: UserState = UserState.Unknown;

  user: Writable<AppUser> = writable<AppUser>(this.localUser);
  userState: Writable<UserState> = writable<UserState>(this.localUserState);
  appEvents: Writable<AppEvent> = writable<AppEvent>({
    type: EventType.Initialize,
    additionalInfo: "",
  });

  defaultServer: string = import.meta.env.VITE_CMS_SERVICE;

  constructor() {
    // In case no default server is set, just try to work with localhost
    if (!this.defaultServer) this.defaultServer = "http://localhost:8080";

    if (browser) {
      this.auth.onAuthStateChanged((u: User) => {
        // if u is undefined, means we don't know user state
        // if u is null, means user is signed out
        // if u is an object, means user is signed in
        this.currentUser = u;

        if (!u) {
          this.localUserState = UserState.SignedOut;
          this.userState.update((n) => this.localUserState);
          //userSignedIn.update((n) => localUserSignedIn);

          if (browser) localStorage.setItem("UserSignedIn", "false");

          this.localUser = undefined;
          this.user.update((n) => this.localUser);

          //if (browser)
          //goto("/home");
        } else {
          console.log("User changed event, user is there.");
          this.localUser = u as AppUser;
          this.user.update((n) => this.localUser);

          this.localUserState = UserState.SignedIn;
          this.userState.update((n) => this.localUserState);

          if (browser) localStorage.setItem("UserSignedIn", "true");

          if (window.location.pathname != "/") {
            // Invalidate all other paths to reload with user, let root forward to /home
            console.log(`invalidating path ${window.location.pathname}`);
            invalidate(window.location.pathname);
          }
        }
      });
    }
  }

  GetUserState(): UserState {
    return this.localUserState;
  }

  GetServer(type: string): string {
    return this.defaultServer;
  }

  Navigate(path: string): void {
    if (browser) goto(path);
  }

  GetIdToken(): Promise<string> {
    return new Promise<string>((resolve, reject) => {
      if (this.currentUser) {
        this.currentUser
          .getIdToken(/* forceRefresh */ true)
          .then(function (idToken) {
            resolve(idToken);
          });
      } else {
        console.log("No user found, sending empty token..");
        resolve("");
      }
    });
  }

  SignIn(provider: SignInProvider): void {
    if (provider === SignInProvider.Google) this.signInWithGoogle();
    else {
      // TODO error handling
      console.error(`Unknown error provider sent ${provider}`);
    }
  }

  signInWithGoogle(): void {
    const auth = getAuth();
    signInWithRedirect(auth, this.googleProvider);
  }

  SignOut(): void {
    signOut(this.auth)
      .then(() => {
        console.log("SignOut successful");
      })
      .catch((error) => {
        // An error happened.
        console.error(error);
      });
  }

  GetPosts(start: number, limit: number): Promise<PostOverview[]> {
    return new Promise<PostOverview[]>((resolve, reject) => {
      this.GetIdToken().then((idToken) => {
        fetch(this.defaultServer + `/posts?start=${start}&limit=${limit}`, {
          method: "GET",
          headers: {
            Accept: "application/json",
            Authorization: "Bearer " + idToken,
          },
        })
          .then((response) => {
            return response.json();
          })
          .then((data: PostOverview[]) => {
            resolve(data);
          });
      });
    });
  }

  GetTaggedPosts(
    tagName: string,
    start: number,
    limit: number
  ): Promise<PostOverview[]> {
    return new Promise<PostOverview[]>((resolve, reject) => {
      this.GetIdToken().then((idToken) => {
        fetch(
          this.defaultServer + `/tags/${tagName}?start=${start}&limit=${limit}`,
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
          .then((data: PostOverview[]) => {
            resolve(data);
          });
      });
    });
  }

  GetPopularPosts(): Promise<PostOverview[]> {
    return new Promise<PostOverview[]>((resolve, reject) => {
      this.GetIdToken().then((idToken) => {
        fetch(this.defaultServer + "/posts/popular?limit=5", {
          method: "GET",
          headers: {
            Accept: "application/json",
            Authorization: "Bearer " + idToken,
          },
        })
          .then((response) => {
            return response.json();
          })
          .then((data: PostOverview[]) => {
            resolve(data);
          });
      });
    });
  }

  GetPost(postId: string): Promise<Post> {
    return new Promise<Post>((resolve, reject) => {
      this.GetIdToken().then((idToken) => {
        fetch(this.defaultServer + "/posts/" + postId, {
          method: "GET",
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

  CreatePost(postData: FormData): Promise<Post> {
    return new Promise<Post>((resolve, reject) => {
      this.GetIdToken()
        .then((idToken) => {
          fetch(this.defaultServer + "/posts", {
            body: postData,
            method: "POST",
            headers: {
              Accept: "application/json",
              Authorization: "Bearer " + idToken,
            },
          })
            .then((response) => {
              return response.json();
            })
            .then((data) => {
              resolve(data);
            });
        })
        .catch((error) => {
          console.error(error);
        });
    });
  }

  UpdatePost(postId: string, postData: FormData): Promise<Post> {
    return new Promise<Post>((resolve, reject) => {
      this.GetIdToken()
        .then((idToken) => {
          fetch(this.defaultServer + "/posts/" + postId, {
            body: postData,
            method: "PUT",
            headers: {
              Accept: "application/json",
              Authorization: "Bearer " + idToken,
            },
          })
            .then((response) => {
              return response.json();
            })
            .then((data) => {
              this.Navigate("/");
            });
        })
        .catch((error) => {
          console.error(error);
        });
    });
  }

  DeletePost(postId: string): Promise<boolean> {
    return new Promise<boolean>((resolve, reject) => {
      this.GetIdToken()
        .then((idToken) => {
          fetch(this.defaultServer + "/posts/" + postId, {
            method: "DELETE",
            headers: {
              Accept: "application/json",
              Authorization: "Bearer " + idToken,
            },
          }).then((response) => {
            resolve(true);
          });
        })
        .catch((error) => {
          console.error(error);
          resolve(false);
        });
    });
  }

  SendEvent(eventType: EventType, additionalInfo?: string) {
    let newEvent: AppEvent = {
      type: eventType,
      additionalInfo: additionalInfo,
    };

    this.appEvents.update((n) => newEvent);
  }

  UpVote(postId: string, user: AppUser): Promise<PostOverview> {
    return new Promise((resolve, reject) => {
      this.GetIdToken().then((idToken) => {
        fetch(this.defaultServer + "/posts/" + postId + "/upvote", {
          method: "POST",
          headers: {
            Accept: "application/json",
            Authorization: "Bearer " + idToken,
          },
        })
          .then((response) => {
            return response.json();
          })
          .then((data: PostOverview) => {
            resolve(data);
          });
      });
    });
  }

  RemoveUpVote(postId: string, user: AppUser): Promise<PostOverview> {
    throw new Error("Method not implemented.");
  }

  GetIfUserUpvoted(postId: string, user: AppUser): Promise<boolean> {
    throw new Error("Method not implemented.");
  }

  GetComments(postId: string): Promise<PostComment[]> {
    return new Promise((resolve, reject) => {
      this.GetIdToken().then((idToken) => {
        fetch(this.defaultServer + "/posts/" + postId + "/comments", {
          method: "GET",
          headers: {
            Accept: "application/json",
            Authorization: "Bearer " + idToken,
          },
        })
          .then((response) => {
            return response.json();
          })
          .then((data: PostComment[]) => {
            resolve(data);
          });
      });
    });
  }

  CreateComment(
    postId: string,
    newCommentData: FormData
  ): Promise<PostComment[]> {
    return new Promise((resolve, reject) => {
      this.GetIdToken()
        .then((idToken) => {
          fetch(this.defaultServer + "/posts/" + postId + "/comments", {
            body: newCommentData,
            method: "post",
            headers: {
              Accept: "application/json",
              Authorization: "Bearer " + idToken,
            },
          })
            .then((response) => {
              return response.json();
            })
            .then((data: PostComment[]) => {
              resolve(data);
            });
        })
        .catch((error) => {
          console.error(error);
        });
    });
  }

  RemoveComment(postId: string, commentId: string, user: AppUser): void {
    throw new Error("Method not implemented.");
  }

  UpvoteComment(postId: string, commentId: string): Promise<PostComment> {
    return new Promise<PostComment>((resolve, reject) => {
      this.GetIdToken().then((idToken) => {
        fetch(
          this.defaultServer +
            "/posts/" +
            postId +
            "/comments/" +
            commentId +
            "/upvote",
          {
            method: "POST",
            headers: {
              Accept: "application/json",
              Authorization: "Bearer " + idToken,
            },
          }
        )
          .then((response) => {
            return response.json();
          })
          .then((data: PostComment) => {
            resolve(data);
          });
      });
    });
  }

  SearchPosts(input: string): Promise<SearchResult[]> {
    return new Promise<SearchResult[]>((resolve, reject) => {
      if (input) {
        this.GetIdToken().then((idToken) => {
          fetch(this.defaultServer + "/posts/search?q=" + input, {
            method: "GET",
            headers: {
              Accept: "application/json",
              Authorization: "Bearer " + idToken,
            },
          })
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

  SearchTags(input: string): Promise<SearchResult[]> {
    return new Promise<SearchResult[]>((resolve, reject) => {
      if (input) {
        this.GetIdToken().then((idToken) => {
          fetch(this.defaultServer + "/tags/search?q=" + input, {
            method: "GET",
            headers: {
              Accept: "application/json",
              Authorization: "Bearer " + idToken,
            },
          })
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

  GetMetadata(): Promise<Metadata> {
    return new Promise<Metadata>((resolve, reject) => {
      this.GetIdToken().then((idToken) => {
        fetch(this.defaultServer + "/admin/data", {
          method: "GET",
          headers: {
            Accept: "application/json",
            Authorization: "Bearer " + idToken,
          },
        })
          .then((response) => {
            return response.json();
          })
          .then((data: Metadata) => {
            resolve(data);
          });
      });
    });
  }

  addComment(
    comments: PostComment[],
    parentCommentId: string,
    newComment: PostComment
  ): boolean {
    if (!parentCommentId) {
      comments.push(newComment);
      return true;
    } else {
      for (let i = 0; i < comments.length; i++) {
        if (comments[i].id === parentCommentId) {
          comments[i].children.push(newComment);
          return true;
        } else if (comments[i].children.length > 0) {
          let result: boolean = this.addComment(
            comments[i].children,
            parentCommentId,
            newComment
          );

          if (result) return true;
        }
      }
    }
  }

  primitiveToBoolean(
    value: string | number | boolean | null | undefined
  ): boolean {
    if (typeof value === "string") {
      return value.toLowerCase() === "true" || !!+value; // here we parse to number first
    }

    return !!value;
  }

  makeid(length) {
    let result = "";
    const characters =
      "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
    const charactersLength = characters.length;
    let counter = 0;
    while (counter < length) {
      result += characters.charAt(Math.floor(Math.random() * charactersLength));
      counter += 1;
    }
    return result;
  }

  // https://gist.github.com/briancavalier/842626?permalink_comment_id=4230644#gistcomment-4230644
  retry = (
    fn: Function,
    retriesLeft = 5,
    interval = 1000,
    intervalMultiplier: (interval: number) => number = (i) => i
  ) =>
    new Promise((resolve, reject) => {
      console.log(
        `Retries left: ${retriesLeft} - Next retry interval: ${interval}`
      );
      fn()
        .then(resolve)
        .catch((error: unknown) => {
          if (retriesLeft === 0) {
            // reject('maximum retries exceeded');
            reject(error);
            return;
          }
          setTimeout(() => {
            // Passing on "reject" is the important part
            this.retry(
              fn,
              retriesLeft - 1,
              intervalMultiplier(interval),
              intervalMultiplier
            ).then(resolve, reject);
          }, interval);
        });
    });
}
