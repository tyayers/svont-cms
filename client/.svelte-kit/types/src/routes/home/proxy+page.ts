// @ts-nocheck
import type { PageLoad } from "./$types";
import { browser } from "$app/environment";
import { UserState } from "../../lib/DataInterface";
import { appService } from "../../lib/DataService";

export const load = (({ params }) => {
  return {
    posts: appService.GetPosts(0, 5),
    popular: appService.GetPopularPosts(),
  };
}) satisfies PageLoad;
;null as any as PageLoad;