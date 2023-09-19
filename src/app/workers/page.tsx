import { IWorker } from "../interfaces/IWorker";
import WorkersTable from "./table";
import { api } from "../_consumer";
import ExperimentSubmission from "./submition";
import { AxiosError } from "axios";

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
            <ExperimentSubmission />
            <WorkersTable workers={workers}/>
        </>
    );
};

export default Homepage;