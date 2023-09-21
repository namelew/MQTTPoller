import { IWorker } from "interfaces/IWorker";
import { api } from "consumers/server";
import WorkersTable from "./table";

const Homepage = async () => {
    let workers:IWorker[] | undefined = undefined;

    try{
        const response = await api.get<IWorker[]>("/worker");
        if (response.status !== 200) {
            alert(response.statusText);
        } else {
            workers = response.data;
        }
    } catch (error) {
        console.log(error);
    }

    return (
        <>
            <WorkersTable workers={workers}/>
        </>
    );
};

export default Homepage;