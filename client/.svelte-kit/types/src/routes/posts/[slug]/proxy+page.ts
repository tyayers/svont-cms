// @ts-nocheck
import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";
import { appService } from "../../../lib/DataService";
import { UserState } from "$lib/DataInterface";

export const load = (({ params }) => {
  return {
    post: appService.GetPost(params.slug),
    comments: appService.GetComments(params.slug),
    popular: appService.GetPopularPosts(),
  };
}) satisfies PageLoad;
;null as any as PageLoad;