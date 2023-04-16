import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";
import { appService } from "../../../lib/DataService";
import { UserState } from "$lib/DataInterface";

export const load: PageLoad = (({ params }) => {
  console.log("enter post load for " + params.slug);

  return {
    post: appService.GetPost(params.slug),
    comments: appService.GetComments(params.slug),
  };
}) satisfies PageLoad;
