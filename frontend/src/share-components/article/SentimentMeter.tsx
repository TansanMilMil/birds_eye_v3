import { Box, Tooltip, Typography } from "@mui/material";
import ThumbUpOffAltIcon from "@mui/icons-material/ThumbUpOffAlt";
import ThumbDownOffAltIcon from "@mui/icons-material/ThumbDownOffAlt";
import { ReactionSentiment } from "../../types/news";

const MIN_SENTIMENT_SAMPLES = 3;

type Props = {
  sentiment?: ReactionSentiment;
};

const toPercent = (count: number, total: number): number =>
  Math.round((count / total) * 100);

export function SentimentMeter({ sentiment }: Props) {
  if (!sentiment) return null;

  const { positive, neutral, negative } = sentiment;
  const total = positive + neutral + negative;
  if (total < MIN_SENTIMENT_SAMPLES) return null;

  const segments = [
    { key: "positive", count: positive, color: "success.main" },
    { key: "neutral", count: neutral, color: "grey.400" },
    { key: "negative", count: negative, color: "error.main" },
  ];

  return (
    <Tooltip
      title={`賛 ${positive}件 / 中立 ${neutral}件 / 否 ${negative}件`}
      placement="top"
    >
      <Box sx={{ marginTop: "0.6rem" }}>
        <Box
          sx={{
            display: "flex",
            justifyContent: "space-between",
            fontSize: "0.75rem",
          }}
        >
          <Typography
            variant="caption"
            sx={{ color: "success.main", display: "inline-flex", alignItems: "center", gap: "0.2rem" }}
          >
            <ThumbUpOffAltIcon aria-label="賛成" sx={{ fontSize: "1rem" }} />
            {toPercent(positive, total)}%
          </Typography>
          <Typography
            variant="caption"
            sx={{ color: "error.main", display: "inline-flex", alignItems: "center", gap: "0.2rem" }}
          >
            <ThumbDownOffAltIcon aria-label="否定" sx={{ fontSize: "1rem" }} />
            {toPercent(negative, total)}%
          </Typography>
        </Box>
        <Box
          role="meter"
          aria-label="反応の賛否"
          aria-valuemin={0}
          aria-valuemax={100}
          aria-valuenow={toPercent(positive, total)}
          sx={{
            display: "flex",
            height: "0.4rem",
            borderRadius: "0.2rem",
            overflow: "hidden",
          }}
        >
          {segments.map(
            (s) =>
              s.count > 0 && (
                <Box key={s.key} sx={{ flexGrow: s.count, bgcolor: s.color }} />
              )
          )}
        </Box>
      </Box>
    </Tooltip>
  );
}
