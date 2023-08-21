import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";
import { appService } from "../../../../lib/DataService";

export const load: PageLoad = (({ params }) => {
  console.log("enter post load for " + params.slug);

  return appService.GetPost(params.slug, true);
}) satisfies PageLoad;
