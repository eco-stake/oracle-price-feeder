package oracle

import (
	"testing"

	"price-feeder/oracle/provider"
	"price-feeder/oracle/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestSuccessFilterTickerDeviations(t *testing.T) {
	providerTickers := make(provider.AggregatedProviderPrices, 4)
	pair := types.CurrencyPair{
		Base:  "ATOM",
		Quote: "USDT",
	}

	atomPrice := sdk.MustNewDecFromStr("29.93")
	atomVolume := sdk.MustNewDecFromStr("1994674.34000000")

	atomTickerPrice := types.TickerPrice{
		Price:  atomPrice,
		Volume: atomVolume,
	}

	providerTickers[provider.ProviderBinance] = map[string]types.TickerPrice{
		pair.Base: atomTickerPrice,
	}
	providerTickers[provider.ProviderHuobi] = map[string]types.TickerPrice{
		pair.Base: atomTickerPrice,
	}
	providerTickers[provider.ProviderKraken] = map[string]types.TickerPrice{
		pair.Base: atomTickerPrice,
	}
	providerTickers[provider.ProviderCoinbase] = map[string]types.TickerPrice{
		pair.Base: {
			Price:  sdk.MustNewDecFromStr("27.1"),
			Volume: atomVolume,
		},
	}

	pricesFiltered, err := FilterTickerDeviations(
		zerolog.Nop(),
		providerTickers,
		make(map[string]sdk.Dec),
	)

	_, ok := pricesFiltered[provider.ProviderCoinbase]
	require.NoError(t, err, "It should successfully filter out the provider using tickers")
	require.False(t, ok, "The filtered ticker deviation price at coinbase should be empty")

	customDeviations := make(map[string]sdk.Dec, 1)
	customDeviations[pair.Base] = sdk.NewDec(2)

	pricesFilteredCustom, err := FilterTickerDeviations(
		zerolog.Nop(),
		providerTickers,
		customDeviations,
	)

	_, ok = pricesFilteredCustom[provider.ProviderCoinbase]
	require.NoError(t, err, "It should successfully not filter out coinbase")
	require.True(t, ok, "The filtered candle deviation price of coinbase should remain")
}

func TestSuccessFilterTickerDeviations2(t *testing.T) {
	pair := types.CurrencyPair{
		Base:  "ATOM",
		Quote: "USDT",
	}

	atomPrice := sdk.MustNewDecFromStr("29.93")
	atomVolume := sdk.MustNewDecFromStr("1994674.34000000")

	atomTickerPrice := types.TickerPrice{
		Price:  atomPrice,
		Volume: atomVolume,
	}

	tickerPrices := map[provider.Name]types.TickerPrice{
		provider.ProviderBinance: atomTickerPrice,
		provider.ProviderHuobi:   atomTickerPrice,
		provider.ProviderKraken:  atomTickerPrice,
		provider.ProviderCoinbase: {
			Price:  sdk.MustNewDecFromStr("27.1"),
			Volume: atomVolume,
		},
	}

	filteredPrices, err := FilterTickerDeviations2(
		zerolog.Nop(),
		pair.String(),
		tickerPrices,
		sdk.NewDec(1),
	)

	require.NoError(t, err)
	require.Len(t, filteredPrices, 3)

	for providerName, tickerPrice := range tickerPrices {
		filteredPrice, found := filteredPrices[providerName]
		if providerName == provider.ProviderCoinbase {
			require.False(t, found)
			continue
		}
		require.True(t, found)
		require.Equal(t, tickerPrice, filteredPrice)
	}

}
