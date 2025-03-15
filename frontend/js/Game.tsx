import React, { createContext, useEffect, useRef, useState } from "react";
import { GameState } from "./types";
import Waiting from "./Waiting";
import Make from "./Make";
import Center from "./Center";
import Trade from "./Trade";
import Leaderboard from "./Leaderboard";

export const StateContext = createContext<GameState | undefined>(undefined);

const Game = () => {
  const eventSourceRef = useRef<EventSource>(undefined);
  const [state, setState] = useState<GameState | undefined>(undefined);

  useEffect(() => {
    if (eventSourceRef.current !== undefined) {
      return;
    }

    const eventSource = new EventSource("http://localhost:8080/event");
    eventSourceRef.current = eventSource;

    eventSource.onmessage = (event) => {
      console.log(event.data);
      setState(JSON.parse(event.data) as GameState);
    };

    eventSource.onerror = () => {
      eventSource.close();
      eventSourceRef.current = undefined;
    };

    return () => {
      eventSource.close();
      eventSourceRef.current = undefined;
    };
  }, []);

  if (state === undefined) {
    return <p>Loading...</p>;
  }

  const renderCurrentView = () => {
    switch (state.view) {
      case "wait":
        return <Waiting />;
      case "make":
        return <Make />;
      case "center":
        return <Center />;
      case "trade":
        return <Trade />;
    }
  };

  return (
    <StateContext.Provider value={state}>
      <div className="flex flex-col items-center justify-center min-h-screen bg-gray-100 space-y-6">
        {renderCurrentView()}
        <Leaderboard />
      </div>
    </StateContext.Provider>
  );
};

export default Game;
