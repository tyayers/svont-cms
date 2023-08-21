// @ts-nocheck
import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";
import { appService } from "../../../../lib/DataService";

export const load = (({ params }) => {
  console.log("enter post load for " + params.slug);

  return appService.GetPost(params.slug, true);
}) satisfies PageLoad;
;null as any as PageLoad;