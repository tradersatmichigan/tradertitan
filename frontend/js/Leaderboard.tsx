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
    { rank: 6, name: "Player 5" },
    { rank: 6, name: "Player 5" },
    { rank: 6, name: "Player 5" },
    { rank: 6, name: "Player 5" },
  ];

  return (
    <div className="bg-white shadow-lg rounded-2xl p-6 max-w-sm w-full">
      <h3 className="text-xl font-semibold text-center text-gray-800">
        Room Leaderboard
      </h3>

      <ul className="mt-4 space-y-2 text-gray-700">
        {data.map(({ rank, name }) => (
          <li
            key={name}
            className="flex justify-between px-4 py-2 bg-gray-100 rounded-lg"
          >
            <span className="font-medium">{rank}.</span>
            <span>{name}</span>
          </li>
        ))}
      </ul>
    </div>
  );
};

export default Leaderboard;
