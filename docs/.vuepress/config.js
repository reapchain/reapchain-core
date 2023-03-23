module.exports = {
  theme: 'cosmos',
  title: 'ReapchainCore Core',
  // locales: {
  //   "/": {
  //     lang: "en-US"
  //   },
  //   "/ru/": {
  //     lang: "ru"
  //   }
  // },
  base: process.env.VUEPRESS_BASE,
  themeConfig: {
    repo: 'reapchain/reapchain-core',
    docsRepo: 'reapchain/reapchain-core',
    docsDir: 'docs',
    editLinks: true,
    label: 'core',
    algolia: {
      id: "BH4D9OD16A",
      key: "59f0e2deb984aa9cdf2b3a5fd24ac501",
      index: "reapchain-core"
    },
    versions: [
      {
        "label": "v0.33",
        "key": "v0.33"
      },
      {
        "label": "v0.34",
        "key": "v0.34"
      },
      {
        "label": "v0.35",
        "key": "v0.35"
      }
    ],
    topbar: {
      banner: false,
    },
    sidebar: {
      auto: true,
      nav: [
        {
          title: 'Resources',
          children: [
            {
              // TODO(creachadair): Figure out how to make this per-branch.
              // See: https://github.com/reapchain/reapchain-core/issues/7908
              title: 'RPC',
              path: 'https://docs.reapchain-core.com/v0.35/rpc/',
              static: true
            },
          ]
        }
      ]
    },
    gutter: {
      title: 'Help & Support',
      editLink: true,
      forum: {
        title: 'ReapchainCore Forum',
        text: 'Join the ReapchainCore forum to learn more',
        url: 'https://forum.cosmos.network/c/reapchain-core',
        bg: '#0B7E0B',
        logo: 'reapchain-core'
      },
      github: {
        title: 'Found an Issue?',
        text: 'Help us improve this page by suggesting edits on GitHub.'
      }
    },
    footer: {
      question: {
        text: 'Chat with ReapchainCore developers in <a href=\'https://discord.gg/vcExX9T\' target=\'_blank\'>Discord</a> or reach out on the <a href=\'https://forum.cosmos.network/c/reapchain-core\' target=\'_blank\'>ReapchainCore Forum</a> to learn more.'
      },
      logo: '/logo-bw.svg',
      textLink: {
        text: 'reapchain-core.com',
        url: 'https://reapchain-core.com'
      },
      services: [
        {
          service: 'medium',
          url: 'https://medium.com/@reapchain-core'
        },
        {
          service: 'twitter',
          url: 'https://twitter.com/reapchain-core_team'
        },
        {
          service: 'linkedin',
          url: 'https://www.linkedin.com/company/reapchain-core/'
        },
        {
          service: 'reddit',
          url: 'https://reddit.com/r/cosmosnetwork'
        },
        {
          service: 'telegram',
          url: 'https://t.me/cosmosproject'
        },
        {
          service: 'youtube',
          url: 'https://www.youtube.com/c/CosmosProject'
        }
      ],
      smallprint:
        'The development of ReapchainCore Core is led primarily by [Interchain GmbH](https://interchain.berlin/). Funding for this development comes primarily from the Interchain Foundation, a Swiss non-profit. The ReapchainCore trademark is owned by ReapchainCore Inc, the for-profit entity that also maintains this website.',
      links: [
        {
          title: 'Documentation',
          children: [
            {
              title: 'Cosmos SDK',
              url: 'https://docs.cosmos.network'
            },
            {
              title: 'Cosmos Hub',
              url: 'https://hub.cosmos.network'
            }
          ]
        },
        {
          title: 'Community',
          children: [
            {
              title: 'ReapchainCore blog',
              url: 'https://medium.com/@reapchain-core'
            },
            {
              title: 'Forum',
              url: 'https://forum.cosmos.network/c/reapchain-core'
            }
          ]
        },
        {
          title: 'Contributing',
          children: [
            {
              title: 'Contributing to the docs',
              url: 'https://github.com/reapchain/reapchain-core'
            },
            {
              title: 'Source code on GitHub',
              url: 'https://github.com/reapchain/reapchain-core'
            },
            {
              title: 'Careers at ReapchainCore',
              url: 'https://reapchain-core.com/careers'
            }
          ]
        }
      ]
    }
  },
  plugins: [
    [
      '@vuepress/google-analytics',
      {
        ga: 'UA-51029217-11'
      }
    ],
    [
      '@vuepress/plugin-html-redirect',
      {
        countdown: 0
      }
    ]
  ]
};
