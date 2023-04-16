import { LitElement } from 'lit';
/**
 * An example element.
 *
 * @slot - This element has a slot
 * @csspart button - The button
 */
export declare class Posts extends LitElement {
    /**
     * Copy for the read the docs hint.
     */
    docsHint: string;
    /**
     * The number of times the button has been clicked.
     */
    count: number;
    posts: PostOverviewCollection;
    connectedCallback(): void;
    render(): import("lit-html").TemplateResult<1>;
    static styles: import("lit").CSSResult;
}
declare global {
    interface HTMLElementTagNameMap {
        'posts-overview': Posts;
    }
}
type PostOverviewCollection = {
    [id: string]: PostOverview;
};
type PostOverview = {
    id: string;
    title: string;
    author: string;
    upvotes: number;
    created: string;
    updated: string;
    fileCount: number;
};
export {};
