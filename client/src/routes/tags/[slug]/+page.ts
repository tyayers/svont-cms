import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";
import { appService } from "../../../lib/DataService";
import { UserState } from "$lib/DataInterface";

export const load: PageLoad = (({ params }) => {
  return {
    tagName: params.slug,
    posts: appService.GetTaggedPosts(params.slug, 0, 5),
    popular: appService.GetPopularPosts(),
  };
}) satisfies PageLoad;
