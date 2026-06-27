import { News } from "../../types/news";
import { Box, Chip, Link, Zoom } from "@mui/material";
import styles from "./MainArticle.module.css";
import { ReactionArea } from "./ReactionArea";

type Props = {
  news: News;
  isDisplayReactions?: boolean;
};

export function MainArticle({ news, isDisplayReactions = false }: Props) {
  const clickSourceBy = (e: React.MouseEvent<HTMLElement, MouseEvent>) => {
    e.preventDefault();
    window.open(news.scrapedUrl, "_blank");
  };

  const transitionDelay = () => Math.floor(Math.random() * 400) + "ms";

  return (
    <div>
      <Zoom in={true} style={{ transitionDelay: transitionDelay() }}>
        <article className={styles.article}>
          {news.articleImageUrl && (
            <a href={news.articleUrl} target="_blank" rel="noreferrer">
              <div className={styles.image}>
                <img src={news.articleImageUrl} alt="" />
              </div>
            </a>
          )}
          <div className={styles.title}>
            <Link
              href={news.articleUrl}
              target="_blank"
              rel="noreferrer"
              sx={{ color: "primary.dark" }}
            >
              {news.title}
            </Link>
          </div>
          <div className={styles.description}>{news.summarizedText}</div>
          <Box sx={{ color: "secondary.main" }}>
            <div className={styles.scrapedDateTime}>{news.scrapedDateTime}</div>
          </Box>
          {news.sourceBy && (
            <Chip label={news.sourceBy} size="small" onClick={clickSourceBy} />
          )}
          {isDisplayReactions && (
            <ReactionArea
              news={news}
              reactionCount={news.reactionCount}
            ></ReactionArea>
          )}
        </article>
      </Zoom>
    </div>
  );
}
