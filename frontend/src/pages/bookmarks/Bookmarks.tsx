import { Box, Typography } from "@mui/material";
import { Masonry } from "@mui/lab";
import { MainArticle } from "../../share-components/article/MainArticle";
import { useBookmarks } from "../../hooks/useBookmarks";

export function Bookmarks() {
  const { bookmarks } = useBookmarks();

  return (
    <div>
      <Box sx={{ maxWidth: "1200px", margin: "0 auto" }}>
        {bookmarks.length === 0 ? (
          <Typography
            sx={{ textAlign: "center", margin: "3rem", color: "common.white" }}
          >
            ブックマークした記事はありません
          </Typography>
        ) : (
          <Masonry
            columns={{ xs: 1, sm: 2, md: 3 }}
            spacing={{ xs: 1, sm: 1, md: 1 }}
          >
            {bookmarks.map((news) => (
              <MainArticle
                key={news.articleUrl}
                news={news}
                isDisplayReactions={true}
              ></MainArticle>
            ))}
          </Masonry>
        )}
      </Box>
    </div>
  );
}
