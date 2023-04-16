export type AppUser = {
  displayName: string;
  email: string;
  phoneNumber: string;
  photoURL: string;
  providerId: string;
  uid: string;
  emailVerified: boolean;
  isAnonymous: boolean;
  refreshToken: string;
}

export type SearchResult = {
  id: string;
  title: string;
};

export type PostOverviewCollection = { 
  [id: string]: PostOverview; 
}

export type PostOverview = {
  id: string;
  title: string;
  summary: string;
  author: string;
  created: string;
  updated: string;
  upvotes: number;
  commentCount: number;
  fileCount: number;
}

export type Post = {
  header: PostOverview;
  content: string;
  files: string[];
}