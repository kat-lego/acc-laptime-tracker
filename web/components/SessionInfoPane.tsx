import { formatMilliseconds, formatUnixToLocalDateTime } from '@/utilities/date'
import { Lap, Session } from '../models/session'
import { getCarName } from '@/utilities/kunosCarMap'

interface SessionInfoPaneParams {
  session: Session
}

export default function SessionInfoPane({ session }: SessionInfoPaneParams) {

  return (
    <div>
      <SessionHero session={session} />
      <SessionLapTable session={session} />
    </div>
  )
}

interface SessionHeroParams {
  session: Session
}

function SessionHero({ session }: SessionHeroParams) {

  return (
    <div className="hero pb-6">
      <div className="hero-content text-center">
        <div className="">
          <h1 className="text-4xl font-bold text-blue-500 uppercase">{session.track}</h1>
          <h2 className="text-3xl font-bold opacity-50 uppercase">{getCarName(session.carModel)}</h2>
          <h2 className="text-3xl font-bold pb-4 opacity-80">{formatUnixToLocalDateTime(session.startTime)}</h2>
          {session.isActive && <div className="badge badge-success">active</div>}
          {!session.isActive && <div className="badge badge-error">inactive</div>}
        </div>
      </div>
    </div>
  )
}


interface SessionLapTableParams {
  session: Session
}

function SessionLapTable({ session }: SessionLapTableParams) {

  return (
    <div className="overflow-x-auto justify-center">
      <table className="table">
        <thead>
          <tr>
            <th>Lap #</th>
            <th>Time</th>
            <th>Time Delta</th>
            <th>Validity</th>
            {Array.from({ length: session.numberOfSectors }, (_, i) => (
              <th key={i + 5}>Sector {i + 1}</th >
            ))}
          </tr>
        </thead>
        <tbody>
          {session.laps.map((lap, _) => (
            <SessionLapTableRow key={lap.lapNumber} lap={lap} />
          ))}
        </tbody>
      </table>
    </div >
  )
}

interface SessionLapTableRowParams {
  lap: Lap
}

function SessionLapTableRow({ lap }: SessionLapTableRowParams) {
  let selectionClass = lap.isActive ? "bg-base-200" : ""
  return (
    <tr className={selectionClass}>
      <th>{lap.lapNumber}</th>
      <td>{formatMilliseconds(lap.lapTime)}</td>
      <td>TODO</td>
      <td>{lap.isValid ? "valid" : "invalid"}</td>
      {lap.lapSectors.map((s, _) => (
        <td key={s.sectorNumber}>{formatMilliseconds(s.sectorTime)}</td>
      ))}
    </tr>
  )
}
