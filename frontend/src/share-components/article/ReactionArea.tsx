import Box from "@mui/material/Box";
import { CircularProgress } from "@mui/material";
import { useState } from "react";
import { BirdsEyeApi } from "../../api/birdsEyeApi";
import { NewsReaction } from "../../types/newsReaction";
import { Reaction } from "./Reaction";
import CommentIcon from "@mui/icons-material/Comment";
import {
  Badge,
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  IconButton,
  Link,
} from "@mui/material";
import styles from "./ReactionArea.module.css";
import { News } from "../../types/news";

type Props = {
  news: News;
  reactionCount: number;
};

export function ReactionArea({ news, reactionCount }: Props) {
  const [reactions, setReactions] = useState<NewsReaction[]>([]);
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [toggle, setToggle] = useState<boolean>(false);

  const getReactions = (e: React.MouseEvent<HTMLElement, MouseEvent>) => {
    e.preventDefault();

    setIsLoading(true);
    setToggle(true);
    BirdsEyeApi.getReactions(news.id)
      .then((result) => {
        setIsLoading(false);
        let reactions = result.reactions;
        reactions = reactions.map((r) => {
          r.scrapedDateTime = new Date(
            Date.parse(r.scrapedDateTime)
          ).toLocaleString();
          return r;
        });
        setReactions(reactions);
      })
      .catch((e) => {
        setIsLoading(false);
      });
  };

  const removeReactions = (e: React.MouseEvent<HTMLElement, MouseEvent>) => {
    e.preventDefault();

    setToggle(false);
  };

  return (
    <div>
      <div>
        {reactionCount >= 1 && (
          <Box sx={{ textAlign: "right" }}>
            <IconButton onClick={getReactions}>
              <Badge badgeContent={reactionCount} color="secondary">
                <CommentIcon />
              </Badge>
            </IconButton>
          </Box>
        )}
      </div>
      {toggle && (
        <div>
          <Dialog
            open={toggle}
            onClose={removeReactions}
            scroll={"paper"}
            aria-labelledby="scroll-dialog-title"
            aria-describedby="scroll-dialog-description"
          >
            <DialogTitle id="scroll-dialog-title">
              <div className={styles.title}>
                <Link href={news.articleUrl} target="_blank" rel="noreferrer">
                  {news.title}
                </Link>
              </div>
            </DialogTitle>
            <DialogContent dividers={true}>
              <div id="scroll-dialog-description" tabIndex={-1}>
                {isLoading && (
                  <Box sx={{ textAlign: "center", margin: "1rem" }}>
                    <CircularProgress color="primary" />
                  </Box>
                )}
                {!isLoading &&
                  [...reactions].map((reaction, i) => (
                    <Reaction
                      key={i}
                      reaction={reaction}
                      newsTitle={news.title}
                      index={i}
                    ></Reaction>
                  ))}
              </div>
            </DialogContent>
            <DialogActions>
              <Button onClick={removeReactions}>Close</Button>
            </DialogActions>
          </Dialog>
        </div>
      )}
    </div>
  );
}
