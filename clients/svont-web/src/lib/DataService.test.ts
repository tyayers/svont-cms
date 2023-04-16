import { goto } from "$app/navigation";
import { browser } from "$app/environment";

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
} from "./DataInterface";
import { writable, type StartStopNotifier, type Writable } from "svelte/store";
import type { AppUser } from "./DataInterface";

// This class implements the DataService interface with pure test data,
// letting the complete application run in a local browser session wo
// a backend.
export class DataServiceTest implements DataService {
  protected localUser: AppUser = undefined;
  protected localUserState: UserState = UserState.Unknown;

  user: Writable<AppUser> = writable<AppUser>(this.localUser);
  userState: Writable<UserState> = writable<UserState>(this.localUserState);
  appEvents: Writable<AppEvent> = writable<AppEvent>({
    type: EventType.Initialize,
    additionalInfo: "",
  });

  protected postOverview: PostOverviewCollection = {
    "1": {
      id: "1",
      title: "I like jelly beans",
      summary: "I like jelly beans, what can I say?",
      authorId: "test@test.com",
      authorDisplayName: "",
      authorProfilePic: "",
      created: "2018-01-01T00:00:00.000+01:00",
      updated: "",
      upvotes: 0,
      commentCount: 0,
      fileCount: 0,
    },
    "2": {
      id: "2",
      title: "I like trees",
      summary: "I like trees, what can I say?",
      authorId: "test@test.com",
      authorDisplayName: "",
      authorProfilePic: "",
      created: "2018-01-01T00:00:00.000+01:00",
      updated: "",
      upvotes: 0,
      commentCount: 0,
      fileCount: 0,
    },
    "3": {
      id: "3",
      title: "I like bicycles",
      summary: "I like bicycles, what can I say?",
      authorId: "test@test.com",
      authorDisplayName: "",
      authorProfilePic: "",
      created: "2018-01-01T00:00:00.000+01:00",
      updated: "",
      upvotes: 0,
      commentCount: 0,
      fileCount: 0,
    },
  };

  protected postContent: { [id: string]: Post } = {
    "1": {
      content: "I like jelly beans, what can I say?",
      files: [],
    },
    "2": {
      content: "I like trees, what can I say?",
      files: [],
    },
    "3": {
      content: "I like trees, what can I say?",
      files: [],
    },
  };

  protected postComments: { [id: string]: PostComment[] } = {
    "1": [
      {
        id: "1",
        created: "2018-01-01T00:00:00.000+01:00",
        updated: "",
        authorId: "1",
        authorDisplayName: "Test User",
        authorProfilePic:
          "https://images.unsplash.com/photo-1494790108377-be9c29b29330?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxzZWFyY2h8MXx8dXNlciUyMHByb2ZpbGV8ZW58MHx8MHx8&w=1000&q=80",
        content: "Great post!",
        upvotes: 0,
        children: [
          {
            id: "2",
            created: "2018-01-01T00:00:00.000+01:00",
            updated: "",
            authorId: "1",
            authorDisplayName: "Test User",
            authorProfilePic:
              "https://images.unsplash.com/photo-1494790108377-be9c29b29330?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxzZWFyY2h8MXx8dXNlciUyMHByb2ZpbGV8ZW58MHx8MHx8&w=1000&q=80",
            content: "I really agree with that!",
            upvotes: 0,
            children: [],
          },
        ],
      },
    ],
  };

  constructor() {
    if (browser) {
      if (
        this.primitiveToBoolean(window.localStorage.getItem("UserSignedIn"))
      ) {
        this.localUserState = UserState.SignedIn;
        console.log(
          `dataservice broadcasting userstate ${this.localUserState}`
        );

        this.userState.set(this.localUserState);
      } else {
        this.localUserState = UserState.SignedOut;
        console.log(
          `dataservice broadcasting userstate ${this.localUserState}`
        );

        this.userState.set(this.localUserState);
      }

      // This is to simulate the delayed user data that is received from firebase auth if a user is signed in.
      setTimeout(() => {
        if (this.localUserState == UserState.SignedIn) this.signInWithGoogle();
      }, 1000);
    }
  }
  UpdatePost(postId: string, postData: FormData): Promise<Post> {
    throw new Error("Method not implemented.");
  }
  UpvoteComment(postId: string, commentId: string): Promise<PostComment> {
    throw new Error("Method not implemented.");
  }

  GetUserState(): UserState {
    return this.localUserState;
  }

  GetIdToken(): Promise<string> {
    throw new Error("Method not implemented.");
  }

  Navigate(path: string): void {
    goto(path);
  }

  SignIn(provider: SignInProvider): void {
    if (provider === SignInProvider.Google) this.signInWithGoogle();
    else {
      // TODO error handling
      console.error(`Unknown error provider sent ${provider}`);
    }
  }

  signInWithGoogle(): void {
    this.localUser = {
      email: "test@test.com",
      displayName: "Test User",
      photoURL:
        "https://images.unsplash.com/photo-1494790108377-be9c29b29330?ixlib=rb-4.0.3&ixid=MnwxMjA3fDB8MHxzZWFyY2h8MXx8dXNlciUyMHByb2ZpbGV8ZW58MHx8MHx8&w=1000&q=80",
    } as AppUser;

    this.user.update((n) => this.localUser);

    this.localUserState = UserState.SignedIn;
    this.userState.update((n) => this.localUserState);

    if (browser) localStorage.setItem("UserSignedIn", "true");

    // if (browser)
    goto("/home");
  }

  SignOut(): void {
    this.localUserState = UserState.SignedOut;
    this.userState.update((n) => this.localUserState);
    //userSignedIn.update((n) => localUserSignedIn);

    if (browser) localStorage.setItem("UserSignedIn", "false");

    this.localUser = undefined;
    this.user.update((n) => this.localUser);

    //if (browser)
    goto("/");
  }

  SearchPosts(input: string): Promise<SearchResult[]> {
    console.log("enter searchposts");
    return new Promise((resolve, reject) => {
      let tempResults: SearchResult[] = [];

      if (input) {
        for (const [key, value] of Object.entries(this.postOverview)) {
          if (value && value.title.includes(input)) {
            tempResults.push(value);
          }
        }
      }

      resolve(tempResults);
    });
  }

  GetPosts(): Promise<PostOverviewCollection> {
    return new Promise<PostOverviewCollection>((resolve, reject) => {
      console.log("in loadposts");
      resolve(this.postOverview);
    });
  }

  GetPost(postId: string): Promise<Post> {
    return new Promise((resolve, reject) => {
      if (this.postContent[postId]) {
        let tempPost = this.postContent[postId];
        tempPost.header = this.postOverview[postId];
        resolve(tempPost);
      } else reject(`Post ${postId} not found!`);
    });
  }

  CreatePost(postData: FormData): Promise<Post> {
    return new Promise((resolve, reject) => {
      var newPost: Post = {
        header: {
          id: this.makeid(8),
          title: postData.get("title").toString(),
          summary: postData.get("summary").toString(),
          authorId: "test@test.com",
          authorDisplayName: "",
          authorProfilePic: "",
          created: "2018-01-01T00:00:00.000+01:00",
          updated: "",
          upvotes: 0,
          commentCount: 0,
          fileCount: 0,
        },
        content: postData.get("content").toString(),
        files: [],
      };

      this.postOverview[newPost.header.id] = newPost.header;
      this.postContent[newPost.header.id] = newPost;

      resolve(newPost);
    });
  }

  DeletePost(postId: string): Promise<boolean> {
    throw new Error("Method not implemented.");
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
      if (this.postOverview[postId]) {
        this.postOverview[postId].upvotes++;
        resolve(this.postOverview[postId]);
      } else reject("Post not found!");
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
      if (this.postComments[postId]) {
        resolve(this.postComments[postId]);
      } else {
        resolve([]);
      }
    });
  }

  CreateComment(
    postId: string,
    newCommentData: FormData
  ): Promise<PostComment[]> {
    return new Promise((resolve, reject) => {
      let parentCommentId: string = "";
      if (newCommentData.has("parentCommentId"))
        parentCommentId = newCommentData.get("parentCommentId").toString();

      let newComment: PostComment = {
        id: this.makeid(8),
        created: new Date().toISOString(),
        updated: "",
        authorId: newCommentData.get("authorId").toString(),
        authorDisplayName: newCommentData.get("authorDisplayName").toString(),
        authorProfilePic: newCommentData.get("authorProfilePic").toString(),
        content: newCommentData.get("content").toString(),
        upvotes: 0,
        children: [],
      };

      if (!this.postComments[postId]) {
        this.postComments[postId] = [newComment];
      } else {
        this.addComment(this.postComments[postId], parentCommentId, newComment);
      }

      resolve(this.postComments[postId]);
    });
  }

  RemoveComment(postId: string, commentId: string, user: AppUser): void {
    throw new Error("Method not implemented.");
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
}
