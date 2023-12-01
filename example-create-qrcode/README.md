Prepare the config

Copy from sample config

    cp config-sample.propses config.props

Change config values

# not used at the moment
APP_NAME=xendit-intgr
APP_DEBUG=false

# give write access to MONEY IN Products
KEY_WRITE_MONEY_IN=xnd_development_xxx

# give read access to MONEY IN Products
KEY_READ_MONEY_IN=xnd_development_xxx

# give read access to MONEY and OUT
KEY_READ_ALL_IN=xnd_development_xxx

# Set invoice duration
INVOICE_DURATION=86400

# Testing Only
SAMPLE_VA_ACCOUNT_ID=xxx