import type { Writable } from "svelte/store";

export interface DataService {
  user: Writable<AppUser>;
  userState: Writable<UserState>;
  appEvents: Writable<AppEvent>;
  defaultServer: string;

  // User
  SignIn(provider: SignInProvider): void;
  SignOut(): void;
  GetIdToken(): Promise<string>;
  GetUserState(): UserState;
  Navigate(path: string): void;

  // Posts
  SearchPosts(input: string): Promise<SearchResult[]>;
  SearchTags(input: string): Promise<SearchResult[]>;
  GetPosts(start: number, limit: number): Promise<PostOverview[]>;
  GetPopularPosts(): Promise<PostOverview[]>;
  GetPost(postId: string): Promise<Post>;
  CreatePost(postData: FormData): Promise<Post>;
  UpdatePost(postId: string, postData: FormData): Promise<Post>;
  DeletePost(postId: string): Promise<boolean>;

  // Upvotes
  UpVote(postId: string, user: AppUser): Promise<PostOverview>;
  RemoveUpVote(postId: string, user: AppUser): Promise<PostOverview>;
  GetIfUserUpvoted(postId: string, user: AppUser): Promise<boolean>;

  // Comments
  GetComments(postId: string): Promise<PostComment[]>;
  CreateComment(
    postId: string,
    newCommentData: FormData
  ): Promise<PostComment[]>;
  UpvoteComment(postId: string, commentId: string): Promise<PostComment>;
  RemoveComment(postId: string, commentId: string, user: AppUser): void;

  // Events
  SendEvent(eventType: EventType, additionalInfo?: string);
}

export type AppUser = {
  email: string;
  uid?: string;
  displayName?: string;
  phoneNumber?: string;
  photoURL?: string;
  providerId?: string;
  emailVerified?: boolean;
  isAnonymous?: boolean;
  refreshToken?: string;
};

export type SearchResult = {
  id: string;
  title: string;
};

export type PostOverviewCollection = {
  [id: string]: PostOverview;
};

export type PostOverview = {
  id: string;
  title: string;
  summary: string;
  authorId: string;
  authorDisplayName: string;
  authorProfilePic: string;
  created: string;
  updated: string;
  upvotes: number;
  tags: string[];
  commentCount: number;
  fileCount: number;
};

export type Post = {
  header?: PostOverview;
  content: string;
  files: string[];
};

export type PostComment = {
  id: string;
  created: string;
  updated: string;
  authorId: string;
  authorDisplayName: string;
  authorProfilePic: string;
  content: string;
  upvotes: number;
  children: PostComment[];
};

export type AppEvent = {
  type: EventType;
  additionalInfo: string;
};

export enum HeaderButton {
  NewPost,
  Submit,
}

export enum UserState {
  Unknown,
  SignedOut,
  SignedIn,
}

export enum SignInProvider {
  Google,
}

export enum EventType {
  Initialize,
  Cancel,
}
