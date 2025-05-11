import Image from "next/image";
import SessionInfoPane from "@/components/SessionInfoPane";
import { Session as AccSession } from "@/models/session";

const sessionData: AccSession[] = [
  {
    id: "20250509T093256Z",
    startTime: 1746886033,
    sessionType: "ACC_HOTLAP",
    track: "nurburgring",
    carModel: "bmw_m4_gt3",
    player: "anonymous",
    numberOfSectors: 3,
    previousLapTime: 2147483647,
    bestLapTime: 2147483647,
    completedLaps: 2,
    isActive: false,
    laps: [
      {
        lapNumber: 1,
        lapTime: 7647,
        isValid: true,
        isActive: false,
        lapSectors: [
          {
            sectorNumber: 1,
            sectorTime: 2500,
            isActive: false
          },
          {
            sectorNumber: 2,
            sectorTime: 2500,
            isActive: false
          },
          {
            sectorNumber: 3,
            sectorTime: 2647,
            isActive: false
          },
        ],
      },
      {
        lapNumber: 2,
        lapTime: 7590,
        isValid: false,
        isActive: true,
        lapSectors: [
          {
            sectorNumber: 1,
            sectorTime: 2490,
            isActive: false
          },
          {
            sectorNumber: 2,
            sectorTime: 2500,
            isActive: false
          },
          {
            sectorNumber: 3,
            sectorTime: 2600,
            isActive: false
          },
        ],
      },
    ],
  },
];

export default function Home() {
  return (
    <div
      className="drawer justify-items-center min-h-screen gap-15
      font-[family-name:var(--font-geist-sans)] pb-16 overflow-x-hidden">

      <input id="my-drawer" type="checkbox" className="drawer-toggle" />

      <div className="drawer-content">

        <header className="sticky top-0 z-50 navbar bg-app-black h-fit">

          <div className="navbar-start">
            <label htmlFor="my-drawer" className="btn btn-primary drawer-button">
              <DropDownIcon />
            </label>
          </div>

          <div className="navbar-end">
            <a className="btn btn-ghost text-xl">
              ACC LAPTIME TRACKER
            </a>
          </div>

        </header>

        <main className="relative flex flex-col gap-8 items-start sm:items-start
          max-w-screen-lg w-full px-6 xs:px-1">
          <SessionInfoPane session={sessionData[0]} />
        </main>

        <footer className="row-start-3 flex gap-[24px] flex-wrap items-center justify-center">
          <a
            className="flex items-center gap-2 hover:underline hover:underline-offset-4"
            href="https://nextjs.org/learn?utm_source=create-next-app&utm_medium=default-template-tw&utm_campaign=create-next-app"
            target="_blank"
            rel="noopener noreferrer"
          >
            <Image
              aria-hidden
              src="/file.svg"
              alt="File icon"
              width={16}
              height={16}
            />
            Learn
          </a>
          <a
            className="flex items-center gap-2 hover:underline hover:underline-offset-4"
            href="https://vercel.com/templates?framework=next.js&utm_source=create-next-app&utm_medium=default-template-tw&utm_campaign=create-next-app"
            target="_blank"
            rel="noopener noreferrer"
          >
            <Image
              aria-hidden
              src="/window.svg"
              alt="Window icon"
              width={16}
              height={16}
            />
            Examples
          </a>
          <a
            className="flex items-center gap-2 hover:underline hover:underline-offset-4"
            href="https://nextjs.org?utm_source=create-next-app&utm_medium=default-template-tw&utm_campaign=create-next-app"
            target="_blank"
            rel="noopener noreferrer"
          >
            <Image
              aria-hidden
              src="/globe.svg"
              alt="Globe icon"
              width={16}
              height={16}
            />
            Go to nextjs.org →
          </a>
        </footer>
      </div>
      <div className="drawer-side">
        <label htmlFor="my-drawer" aria-label="close sidebar" className="drawer-overlay"></label>
        <ul className="list rounded-box shadow-md bg-base-100 menu min-h-full w-80 p-4">
          <li className="list-row">Start Time, Track, Carmodel</li>
        </ul>
      </div>
    </div>
  );
}


function DropDownIcon() {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      className="h-5 w-5"
      fill="none"
      viewBox="0 0 24 24"
      stroke="currentColor">
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth="2"
        d="M4 6h16M4 12h16M4 18h7" />
    </svg>
  )
}
