import { IWorker } from "../interfaces/IWorker";
import WorkersTable from "./table";
import { api } from "../_consumer";

const Homepage = async () => {
    const response = await api.get<IWorker[]>("/worker");
    let workers:IWorker[] | undefined = undefined;

    if (response.status !== 200) {
        alert(response.statusText);
    } else {
        workers = response.data;
    }

    return (
        <>
            <WorkersTable workers={workers}/>
        </>
    );
};

export default Homepage;