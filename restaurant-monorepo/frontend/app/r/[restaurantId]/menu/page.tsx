"use client";

import { useState, useEffect } from 'react';
import api from '@/lib/api';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { useParams } from 'next/navigation';

interface Category {
  ID: string;
  Name: string;
}

interface MenuItem {
  ID: string;
  Name: string;
  Price: number;
  CategoryID: string;
  Description: string;
}

export default function PublicMenuPage() {
  const params = useParams();
  const restaurantId = params.restaurantId as string;
  
  const [categories, setCategories] = useState<Category[]>([]);
  const [items, setItems] = useState<MenuItem[]>([]);

  useEffect(() => {
    const fetchMenu = async () => {
      if (!restaurantId) return;
      try {
        const res = await api.get(`/public/restaurants/${restaurantId}/menu`);
        setCategories(res.data.categories);
        setItems(res.data.items);
      } catch (e) {
        console.error(e);
      }
    };
    fetchMenu();
  }, [restaurantId]);

  return (
    <div className="container mx-auto p-8">
      <h1 className="text-4xl font-bold mb-8 text-center">Menu</h1>
      
      {categories.map(category => (
        <div key={category.ID} className="mb-12">
          <h2 className="text-2xl font-semibold mb-6 border-b pb-2">{category.Name}</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            {items.filter(i => i.CategoryID === category.ID).map(item => (
              <Card key={item.ID} className="hover:shadow-lg transition-shadow">
                <CardHeader>
                  <CardTitle className="flex justify-between">
                    <span>{item.Name}</span>
                    <span>${item.Price}</span>
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <p className="text-gray-600">{item.Description}</p>
                </CardContent>
              </Card>
            ))}
          </div>
        </div>
      ))}
    </div>
  );
}
