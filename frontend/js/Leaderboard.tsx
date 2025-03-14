import React, { useContext } from "react";
import { StateContext } from "./Game";

const Leaderboard = () => {
  const state = useContext(StateContext);

  const data = [
    { rank: 1, name: "Conner" },
    { rank: 2, name: "Player 1" },
    { rank: 2, name: "Player 2" },
    { rank: 4, name: "Player 3" },
    { rank: 5, name: "Player 4" },
    { rank: 6, name: "Player 5" },
  ];

  return (
    <>
      <h3>Room leaderboard</h3>
      <ul className="list-none">
        {data.map(({ rank, name }) => (
          <li key={name}>
            {rank}. {name}
          </li>
        ))}
      </ul>
    </>
  );
};

export default Leaderboard;
