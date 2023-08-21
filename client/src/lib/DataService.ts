import { DataServiceTest } from "./DataService.test";
import { DataServiceGoogle } from "./DataService.google";
import type { DataService } from "./DataInterface";
import logo from "../assets/svelte.png";

export let LogoPath: string = logo;

// Change this to any service implementation that implements the
// DataService interface.
export let appService: DataService = new DataServiceGoogle();
