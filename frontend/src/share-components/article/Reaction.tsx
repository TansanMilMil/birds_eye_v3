import styles from "./Reaction.module.css";
import { NewsReaction } from "../../types/newsReaction";
import reactStringReplace from "react-string-replace";
import { Avatar, Box, Grow, Link } from "@mui/material";
import PersonIcon from "@mui/icons-material/Person";
import TwitterIcon from "@mui/icons-material/Twitter";
import { blue, lightBlue } from "@mui/material/colors";

type Props = {
  reaction: NewsReaction;
  newsTitle: string;
  index: number;
};

export function Reaction({ reaction, newsTitle, index }: Props) {
  const changeTitleColor = (text: string, title: string) => {
    if (!text) return text;

    let dom: React.ReactNode[] = [text];
    dom = reactStringReplace(
      dom,
      new RegExp("(" + title + ")", "g"),
      (match, i) => (
        <Link
          key={match + i}
          href={match}
          target="_blank"
          rel="noreferrer"
          className="ref-title"
        >
          {match}
        </Link>
      )
    );
    dom = reactStringReplace(dom, /(https?:\/\/\S+)/g, (match, i) => (
      <Link key={match + i} href={match} target="_blank" rel="noreferrer">
        {match}
      </Link>
    ));
    dom = reactStringReplace(dom, /#(\S+)/g, (match, i) => (
      <Link
        key={match + i}
        href={`https://twitter.com/hashtag/${match}`}
        target="_blank"
        rel="noreferrer"
      >
        #{match}
      </Link>
    ));
    return dom;
  };

  const transitionDelay = () => index * 100 + "ms";

  return (
    <div>
      <Grow in={true} style={{ transitionDelay: transitionDelay() }}>
        <div className={styles.reactionArea}>
          <div className={styles.author}>
            {(() => {
              switch (reaction.author) {
                case "twitter user":
                  return (
                    <Avatar
                      sx={{ width: 24, height: 24, bgcolor: lightBlue[500] }}
                    >
                      <TwitterIcon />
                    </Avatar>
                  );
                case "hatena user":
                  return (
                    <Avatar sx={{ width: 24, height: 24, bgcolor: blue[500] }}>
                      B!
                    </Avatar>
                  );
                default:
                  return (
                    <Avatar sx={{ width: 24, height: 24 }}>
                      <PersonIcon />
                    </Avatar>
                  );
              }
            })()}
            <Box sx={{ paddingLeft: "0.5rem" }}>
              <Link
                href={reaction.commentUrl}
                target="_blank"
                rel="noreferrer"
                sx={{ color: "primary.dark" }}
              >
                {reaction.author}
              </Link>
            </Box>
          </div>
          <div className={styles.comment}>
            {changeTitleColor(reaction.comment, newsTitle)}
          </div>
          <Box sx={{ color: "secondary.main" }}>
            <div className={styles.scrapedDateTime}>
              {reaction.scrapedDateTime}
            </div>
          </Box>
        </div>
      </Grow>
    </div>
  );
}
