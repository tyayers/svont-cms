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
  GetTaggedPosts(
    tagName: string,
    start: number,
    limit: number
  ): Promise<PostOverview[]>;
  GetPost(postId: string): Promise<Post>;
  CreatePost(postData: FormData): Promise<Post>;
  UpdatePost(postId: string, postData: FormData): Promise<Post>;
  DeletePost(postId: string): Promise<boolean>;
  GetServer(type: string): string;

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

  // Admin
  GetMetadata(): Promise<Metadata>;
  DoRefresh(): Promise<Metadata>;
  DoPersist(): Promise<Metadata>;
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
  count: number;
};

export type PostOverviewCollection = {
  [id: string]: PostOverview;
};

export type PostOverview = {
  id: string;
  title: string;
  summary: string;
  image: string;
  draft: boolean;
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

export type Metadata = {
  postCount: number;
  draftCount: number;
  deletedCount: number;
  userCount: number;
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

export function ToTitleCase(input: string) {
  var smallWords =
    /^(a|an|and|as|at|but|by|en|for|if|in|nor|of|on|or|per|the|to|v.?|vs.?|via)$/i;
  var alphanumericPattern = /([A-Za-z0-9\u00C0-\u00FF])/;
  var wordSeparators = /([ :–—-])/;

  return input
    .split(wordSeparators)
    .map(function (current, index, array) {
      if (
        /* Check for small words */
        current.search(smallWords) > -1 &&
        /* Skip first and last word */
        index !== 0 &&
        index !== array.length - 1 &&
        /* Ignore title end and subtitle start */
        array[index - 3] !== ":" &&
        array[index + 1] !== ":" &&
        /* Ignore small words that start a hyphenated phrase */
        (array[index + 1] !== "-" ||
          (array[index - 1] === "-" && array[index + 1] === "-"))
      ) {
        return current.toLowerCase();
      }

      /* Ignore intentional capitalization */
      if (current.substr(1).search(/[A-Z]|\../) > -1) {
        return current;
      }

      /* Ignore URLs */
      if (array[index + 1] === ":" && array[index + 2] !== "") {
        return current;
      }

      /* Capitalize the first letter */
      return current.replace(alphanumericPattern, function (match) {
        return match.toUpperCase();
      });
    })
    .join("");
}
