// @ts-nocheck
import type { PageLoad } from "./$types";
import { browser } from "$app/environment";
import { UserState } from "../../lib/DataInterface";
import { appService } from "../../lib/DataService";

export const load = (({ params }) => {
  return {
    metadata: appService.GetMetadata(),
  };
}) satisfies PageLoad;
;null as any as PageLoad;