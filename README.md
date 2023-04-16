# SVONT - Svelt Content Collaboration

Svont (**sv**elt+c**ont**ent) is a simple, feature-rich Content Collaboration System (CCS) using Svelte and a serverless / NODB stack (think serverless Wordpress with Svelte). You can deploy Svont either in an all-in-one container, or distributed across static hosting and a backend service.

## Why another CMS/CCS?

Because most CMS's require a database like MySQL or Postgresql, and so are heavier and more resource intensive to run. Svont aims to be extremely _svelte_, like the framework itself, and run very efficiently using cloud serverless platforms and static hosting services.

This project was initially built for the SveltSociety Hackathon 2023 (Feb 2023 - April 2023).

## Why Svelte?

Many frameworks were tested for this project, including solid, react, preact, vue, etc... You can find the tests in the clients/test-clients directory, but long-story-short svelte won out on capabilities, size & speed (preact was the other finalist). A blog article is coming with more detailed results of the comparison.

## Which use-cases is Svont good for?

Svont is ideal for use-cases such as:

1. Team collaboration site, where users can post content updates with document attachments, get likes, leave comments, etc... Think Sharepoint or Confluence, but simpler and more lightweight and less resource intensive to run (but still with the most important features).
2. Public blog with content posts, search, and the possibility for registered users to like and comment on content. Think Wordpress blogs, with both public and private registered user content & interaction.
3. Micro social platforms to quickly and easily host content and users for a specific topic, with social features like likes and comments.

## Currently supported features

These features are currently supported:

- User registration and sign-in (currently using Firebase Auth)
- Content posting, editing and deleting (currently using CKEditor as content editor)
- Content file attachments
- Content likes and comments
- Built-in theme inspired by Medium.com
- Full-text search of content (using the blevee framework, similar to Elasticsearch but embedded/lightweight)

## Roadmap

Here is the roadmap of planned features:

- Post overview pagination - currently it's a single long list of posts, pagniation needed (prio 1)
- RSS feed support - automatically generate and publish RSS feed of content
- Theme support - admins can choose from different themes to use, also with dark mode support
  - Material Design theme
  - Bootstrap (Twitter) theme
  - Others?
- Admin console to manage site configuration
- Email notification support - get notified by email for new comment, likes, follow specific users, etc..
- User statistics and overview - have a view per user of posts by user, comments, other data, possibility to follow, etc..
