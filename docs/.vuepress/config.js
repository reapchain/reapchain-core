module.exports = {
  theme: 'cosmos',
  title: 'Reapchain Core',
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
    repo: 'reapchain/reapchain',
    docsRepo: 'reapchain/reapchain',
    docsDir: 'docs',
    editLinks: true,
    label: 'core',
    algolia: {
      id: "BH4D9OD16A",
      key: "59f0e2deb984aa9cdf2b3a5fd24ac501",
      index: "reapchain"
    },
    versions: [
      {
        "label": "v0.32",
        "key": "v0.32"
      },
      {
        "label": "v0.33",
        "key": "v0.33"
      },
      {
        "label": "v0.34",
        "key": "v0.34"
      },
      {
        "label": "master",
        "key": "master"
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
              title: 'Developer Sessions',
              path: '/DEV_SESSIONS.html'
            },
            {
              title: 'RPC',
              path: 'https://docs.reapchain.com/master/rpc/',
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
        title: 'Reapchain Forum',
        text: 'Join the Reapchain forum to learn more',
        url: 'https://forum.cosmos.network/c/reapchain',
        bg: '#0B7E0B',
        logo: 'reapchain'
      },
      github: {
        title: 'Found an Issue?',
        text: 'Help us improve this page by suggesting edits on GitHub.'
      }
    },
    footer: {
      question: {
        text: 'Chat with Reapchain developers in <a href=\'https://discord.gg/vcExX9T\' target=\'_blank\'>Discord</a> or reach out on the <a href=\'https://forum.cosmos.network/c/reapchain\' target=\'_blank\'>Reapchain Forum</a> to learn more.'
      },
      logo: '/logo-bw.svg',
      textLink: {
        text: 'reapchain.com',
        url: 'https://reapchain.com'
      },
      services: [
        {
          service: 'medium',
          url: 'https://medium.com/@reapchain'
        },
        {
          service: 'twitter',
          url: 'https://twitter.com/reapchain_team'
        },
        {
          service: 'linkedin',
          url: 'https://www.linkedin.com/company/reapchain/'
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
        'The development of Reapchain Core is led primarily by [Interchain GmbH](https://interchain.berlin/). Funding for this development comes primarily from the Interchain Foundation, a Swiss non-profit. The Reapchain trademark is owned by Reapchain Inc, the for-profit entity that also maintains this website.',
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
              title: 'Reapchain blog',
              url: 'https://medium.com/@reapchain'
            },
            {
              title: 'Forum',
              url: 'https://forum.cosmos.network/c/reapchain'
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
              title: 'Careers at Reapchain',
              url: 'https://reapchain.com/careers'
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
    ]
  ]
};
