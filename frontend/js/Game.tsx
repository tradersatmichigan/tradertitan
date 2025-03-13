import React, { useEffect } from "react";

const Game = () => {
  useEffect(() => {
    const eventSource = new EventSource("event");

    eventSource.onmessage = (event) => {
      console.log(JSON.parse(event.data));
    };

    eventSource.onerror = (error) => {
      console.error("SSE Error:", error);
      eventSource.close();
    };

    return () => {
      eventSource.close();
    };
  }, []);

  return <p>This is the game</p>;
};

export default Game;
