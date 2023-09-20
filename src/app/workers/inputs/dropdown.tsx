import { Select, FormLabel, Box} from "@chakra-ui/react";
import { HTMLInputTypeAttribute } from "react";

interface Option {
    key: string,
    name?: string,
    value?: number | string | readonly string[],
}

interface Props {
    label?: string
    name?: string,
    value?: number | string | readonly string[],
    options: Option[],
    onChange?: (event:React.FormEvent<any>) => void
}

const DropDown = ({ label, name, value, onChange, options } : Props) => {
    return (
        <Box minW='3xs'>
            {label && <FormLabel>{label}</FormLabel>}
            <Select name={name} value={value} onClick={onChange}>
                {options.map((option) => (
                    <option key={option.key} value={option.value}>{option.name}</option>
                ))}
            </Select>
        </Box>
    );
};

export default DropDown;