import { api } from "consumer";
import { IExperiment } from "interfaces/IExperiment";
import ExperimentsTable from "./table";

const Homepage = async () => {
    let experiments:IExperiment[] | undefined = undefined;

    try{
        const response = await api.get<IExperiment[]>("/experiment");

        if (response.status !== 200) {
            alert(response.statusText);
        } else {
            experiments = response.data;
        }
    } catch (error) {
        console.log(error);
    }
    
    return (
        <ExperimentsTable experiments={experiments}/>
    );
};

export default Homepage;