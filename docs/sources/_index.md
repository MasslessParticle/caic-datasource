# CAIC datasource for grafana

The CAIC datasource allows you to visualize data from the Colorado Avalanche Information Center about current avalanche conditions in the Colorado backcountry.

## Requirements

There are no specific requirements to run this plugin

## Known limitations

- The plugin only pulls region ratings and aspect angle dangers at this time

## Install the plugin

1. Navigate to [The plugin on Github](https://github.com/MasslessParticle/ciac-datasource).
1. Clone the plugin to your grafana plugins directory
1. Restart Grafana

### Meet compatibility requirements

For this plugin, there are no compatibility requirements.

### Verify that the plugin is installed

1. In Grafana from the left-hand menu, navigate to **Configuration** > **Data sources**.
2. From the top-right corner, click the **Add data source** button.
3. Search for `caic-datasource` in the search field, and hover over the search result.
4. Click the **Select** button for .
   - If you can click the **Select** button, then it is installed.

## Configure the data source

This plugin pulls from the publicly available CAIC website so no specific configuration is needed

## Learn more

- [Colorado Avalanhe Information Center](https://www.avalanche.state.co.us/).
