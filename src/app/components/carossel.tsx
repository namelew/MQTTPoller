'use client'
import React, { useState } from 'react';
import { Box, Button, Center, VStack } from "@chakra-ui/react";

interface Props {
    items: React.ReactNode[]
};

const Carousel = ({ items } : Props) => {
  const [displayIndex, setDisplayIndex] = useState(0);

  const handleNext = () => {
    setDisplayIndex((displayIndex + 1) % items.length);
  };

  const handlePrev = () => {
    setDisplayIndex((displayIndex - 1 + items.length) % items.length);
  };

  return (
    <VStack>
      <Center border="1px" borderColor="gray.200" borderRadius="lg" boxShadow="lg">
        {items.length > 0 ? items[displayIndex] : 'Sem items'}
      </Center>
      {items.length > 1 &&
        <Box>
          <Button onClick={handlePrev} mr="4">Anterior</Button>
          <Button onClick={handleNext}>Próximo</Button>
        </Box>
      }
    </VStack>
  );
};

export default Carousel;
