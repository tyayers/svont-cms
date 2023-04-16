import { LitElement, html } from 'lit';

const tagName = 'posto';

export class PostOverview extends LitElement {
  static properties = {
    name: {},
  };

  constructor() {
    super();
    // Declare reactive properties
    this.name = 'World';
  }

  render() {
    return html` <p>Hello world! From PostOverview</p> `;
  }
}

customElements.define(tagName, PostOverview);